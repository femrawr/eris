<?php

const REDDIT_ACCOUNT_SCHEMA = [
    'user_id' => 'i',
    'title_name' => 'n',
    'description' => 'd',
    'over_18' => ['o', 'bool'],
    'total_karma' => ['t', 'int'],
    'comment_karma' => ['c', 'int'],
    'accept_dms' => ['a', 'bool'],
    'created' => ['m', 'int']
];

const REDDIT_POST_SCHEMA = [
    'user_id' => 'i',
    'post_id' => 'p',
    'sub_reddit' => 's',
    'body' => 'b',
    'title' => 't',
    'upvotes' => ['u', 'int'],
    'vote_ratio' => ['v', 'int'],
    'over_18' => ['o', 'bool'],
    'created' => ['m', 'int']
];

const REDDIT_COMMENT_SCHEMA = [
    'user_id' => 'i',
    'post_id' => 'p',
    'sub_reddit' => 's',
    'parent_id' => 'r',
    'body' => 'b',
    'upvotes' => ['u', 'int'],
    'created' => ['m', 'int']
];