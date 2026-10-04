<?php

declare(strict_types=1);

namespace Habarchy;

/** Non-2xx response from the Habarchy API. */
final class HabarchyException extends \RuntimeException
{
    /** @param array<string,mixed>|null $details */
    public function __construct(
        public readonly int $status,
        public readonly string $errorCode,
        string $message,
        public readonly ?array $details = null,
    ) {
        parent::__construct(sprintf('habarchy: %s (%d): %s', $errorCode, $status, $message), $status);
    }
}
