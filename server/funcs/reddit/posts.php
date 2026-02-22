<?php

/*
    var body - /index
    var method - /index

    con REDDIT_POST_SCHEMA - /schemas/reddit

    fun get_db_file - /utils/database

    fun validate_fields - /utils/validate

    fun num_to_base62 - /utils/encode

    fun compress - /utils/compress
*/

if ($method === 'POST') {
    validate_fields(REDDIT_POST_SCHEMA);

    $db_file = get_db_file('reddit', $body['user_name'], 'posts');
    $saved_at = num_to_base62(time());

    unset($body['type']);
    unset($body['user_name']);

    $compressed = compress($body, REDDIT_POST_SCHEMA);

    file_put_contents(
        $db_file,
        $saved_at . json_encode($compressed) . "\n",
        FILE_APPEND
    );

    echo json_encode(['ok' => true]);
    exit;
}