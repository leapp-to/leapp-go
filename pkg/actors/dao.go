package actors

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// ActorDAO is the interface to the DB Storage operations for the Actor API
type ActorDAO struct {
	DB *sqlx.DB
}

type mustExec interface {
	MustExec(string, ...interface{}) sql.Result
}

type dbGet interface {
	Get(interface{}, string, ...interface{}) error
}

type dbGetMustExec interface {
	mustExec
	dbGet
}

func withTransaction(db *sqlx.DB, txFun func(tx *sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}

	defer func() {
		if paniced := recover(); paniced != nil {
			tx.Rollback()
			panic(paniced) // Panic again
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFun(tx)

	return err
}

func storeHost(db dbGetMustExec, host Host) (hostID int64, err error) {
	db.MustExec(`
		INSERT OR IGNORE INTO host (
			context, hostname
		) VALUES(?, ?)
	`, host.Context, host.Hostname)

	err = db.Get(&hostID, `
		SELECT id FROM host WHERE context = ? AND hostname = ?
	`, host.Context, host.Hostname)

	return hostID, err
}

func storeDataSource(db dbGetMustExec, ds DataSource) (dataSourceID int64, err error) {

	hostID, err := storeHost(db, ds.Host)

	db.MustExec(`
		INSERT OR IGNORE INTO data_source (
			context, host_id, actor, phase
		) VALUES(?, ?, ?, ?)
	`, ds.Context, hostID, ds.Actor, ds.Phase)

	err = db.Get(&dataSourceID, `
		SELECT id
		FROM data_source
		WHERE context = ? AND host_id = ? AND actor = ? AND phase = ?
	`, ds.Context, hostID, ds.Actor, ds.Phase)

	return dataSourceID, err
}

func storeMessage(db *sqlx.Tx, msg *Message) (messageID int64, err error) {
	db.MustExec(`
		INSERT OR IGNORE INTO message_data (
			hash, data
		) VALUES(?, ?)
	`, msg.Message.Hash, msg.Message.Data)

	dataSourceID, err := storeDataSource(db, msg.DataSource)
	db.MustExec(`
		INSERT OR IGNORE INTO message (
			context, stamp, channel, type, data_source_id, message_data_hash
		) VALUES(?, ?, ?, ?, ?, ?)
	`, msg.Context, msg.Stamp, msg.Channel, msg.Type, dataSourceID, msg.Message.Hash)

	err = db.Get(&messageID, `SELECT last_insert_rowid()`)

	return messageID, err
}

// // AddMessage adds a new message which was produced by an actor to the database
// func (dao ActorDAO) AddMessage(msg *Message) error {
// 	return withTransaction(dao.DB, func(tx *sqlx.Tx) (err error) {
// 		_, err = storeMessage(tx, msg)
// 		return
// 	})
// }

// AddAudit adds a new audit entry to the database
func (dao ActorDAO) AddAudit(audit *Audit) error {
	return withTransaction(dao.DB, func(tx *sqlx.Tx) (err error) {

		dataSourceID, err := storeDataSource(tx, audit.DataSource)

		var messageID *int64
		if audit.Message != nil {
			messageID = new(int64)
			*messageID, err = storeMessage(tx, audit.Message)
		}
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO audit (
				event, stamp, context, data_source_id, message_id, data
			)
			VALUES(?, ?, ?, ?, ?, ?)
		`, audit.Event, audit.Stamp, audit.Context, dataSourceID, messageID, audit.Data)

		return
	})
}

/*
	message.id AS id,
	message.context AS context,
	message.stamp AS stamp,
	message.channel AS channel,
	message.type as type,
	data_source.actor as actor,
	data_source.phase as phase,
	msg_data.hash as message_hash,
	msg_data.data as message_data,
	host.hostname as hostname
*/

// GetMessages queries messages from the database based on the input parameters
func (dao ActorDAO) GetMessages(context string, messageTypes []string) (messages []Message, err error) {
	query, args, err := sqlx.In(`
		SELECT
			id, context, stamp, channel, type, actor, phase, message_hash, message_data, hostname
		FROM 
			messages_data
		WHERE context = ? AND type IN (?)
	`, context, messageTypes)
	if err != nil {
		return nil, err
	}

	query = dao.DB.Rebind(query)
	rows, err := dao.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var msg Message
		err = rows.Scan(
			&msg.ID,
			&msg.Context,
			&msg.Stamp,
			&msg.Channel,
			&msg.Type,
			&msg.Actor,
			&msg.Phase,
			&msg.Message.Hash,
			&msg.Message.Data,
			&msg.Hostname)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if messages == nil {
		messages = []Message{}
	}
	return messages, nil
}
