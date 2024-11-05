// guard.ts
export function guard(condition: any, message: string): asserts condition {
    if (!condition) {
        throw new Error(message);
    }
}

export function guardRFC3339(dateString: string | null, message: string): asserts dateString is string {
    if (!dateString) {
        throw new Error(message);
    }

    // RFC3339 regex pattern for basic date validation
    const rfc3339Pattern = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/;

    if (!rfc3339Pattern.test(dateString)) {
        throw new Error(`Invalid date format: ${dateString}. Expected RFC3339 format.`);
    }
}