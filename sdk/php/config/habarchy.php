<?php

return [
    'url' => env('HABARCHY_URL', 'http://localhost:8080'),
    'api_key' => env('HABARCHY_API_KEY', ''),
    // Required when the API key was created with require_signature.
    'sign' => (bool) env('HABARCHY_SIGN', false),
    'timeout' => (int) env('HABARCHY_TIMEOUT', 20),
];
