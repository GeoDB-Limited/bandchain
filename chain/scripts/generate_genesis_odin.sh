#!/bin/bash

# remove old genesis
rm -rf ~/.band*

# initial new node
bandd init node-validator --chain-id odin

# create accounts
echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | bandcli keys add validator1 --recover --keyring-backend test

echo "loyal damage diet label ability huge dad dash mom design method busy notable cash vast nerve congress drip chunk cheese blur stem dawn fatigue" \
    | bandcli keys add validator2 --recover --keyring-backend test

echo "whip desk enemy only canal swear help walnut cannon great arm onion oval doctor twice dish comfort team meat junior blind city mask aware" \
    | bandcli keys add oracle-validator --recover --keyring-backend test

# add accounts to genesis
bandd add-genesis-account validator1 10000000000000loki --keyring-backend test
bandd add-genesis-account validator2 10000000000000loki --keyring-backend test
bandd add-genesis-account oracle-validator 1000000000000000loki --keyring-backend test

# genesis configurations
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

# register initial validators
bandd gentx \
    --amount 100000000loki \
    --node-id 11392b605378063b1c505c0ab123f04bd710d7d7 \
    --pubkey odinvalconspub1addwnpepqgjt0rywpn3s8aqz48lwzxrd02yuz97djzlnw50guxz3w6h9zj3jg06gpvl \
    --name validator1 \
    --details "Alice's Adventures in Wonderland (commonly shortened to Alice in Wonderland) is an 1865 novel written by English author Charles Lutwidge Dodgson under the pseudonym Lewis Carroll." \
    --website "https://www.alice.org/" \
    --ip 172.18.0.11 \
    --keyring-backend test

bandd gentx \
    --amount 100000000loki \
    --node-id 0851086afcd835d5a6fb0ffbf96fcdf74fec742e \
    --pubkey odinvalconspub1addwnpepq0myardtekr2xqnzz87fw3qplqkx4q36rfceuuard5f69de66xwcv5papww \
    --name validator2 \
    --details "Fish is best known for his appearances with Ring of Honor (ROH) from 2013 to 2017, where he wrestled as one-half of the tag team reDRagon and held the ROH World Tag Team Championship three times and the ROH World Television Championship once." \
    --website "https://www.wwe.com/superstars/bobby-fish" \
    --ip 172.18.0.12 \
    --keyring-backend test

bandd gentx \
    --amount 100000000loki \
    --node-id f7343e1aeafb7b20d37e0efdcd331a04528cbf66 \
    --pubkey odinvalconspub1addwnpepqv58qjx2d9yfzhnldw6peaxt6lzv4halc485u2j6z3nlg8t8xph22ylmjef \
    --name oracle-validator \
    --details "Carol Susan Jane Danvers is a fictional superhero appearing in American comic books published by Marvel Comics. Created by writer Roy Thomas and artist Gene Colan." \
    --website "https://www.marvel.com/characters/captain-marvel-carol-danvers" \
    --ip 172.18.0.13 \
    --keyring-backend test

# collect genesis transactions
bandd collect-gentxs

