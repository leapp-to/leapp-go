#!/bin/bash

request_update_doc() {
    git clone "https://$GH_TOKEN@github.com/leapp-to/leapp-to.github.io.git"
    cd leapp-to.github.io
    git config user.name "Travis-CI"
    cp ../index.html.md shins/source/index.html.md
    cd shins
    npm update
    node ./shins.js --minify
    git add -A
    git add pub/js/shins.js
    git commit -m "rebuild pages at $commit"
    git push origin HEAD:master
}

install_npm_deps() {
    npm install -g widdershins
}

convert_documentation() {
    node $NVM_BIN/widdershins -y ./docs/api/api.yaml -o index.html.md
}

#checks if it's merge action
if [[ $TRAVIS_PULL_REQUEST == "false" && $TRAVIS_BRANCH == "master" ]]; then
    changed_files=`git diff --name-only HEAD^`
    commit=$(git rev-parse --short HEAD)
    # checks if merged PR contains any changes in api.yaml
    if [[ $changed_files =~ .*api.yaml ]]; then
        echo "Update documentation has been triggered"
        install_npm_deps
        convert_documentation
        request_update_doc
        exit 0
    fi
fi

echo "Update documentation has not been triggered"
