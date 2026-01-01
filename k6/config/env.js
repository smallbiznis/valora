export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
export const API_KEY = __ENV.API_KEY || '';
export const ORG_ID = __ENV.ORG_ID || '';

export const CUSTOMER_ID = __ENV.CUSTOMER_ID || '';
export const METER_CODE = __ENV.METER_CODE || 'api_calls';

export const ADMIN_SESSION = __ENV.ADMIN_SESSION || '';

export const PAYMENT_PROVIDER = __ENV.PAYMENT_PROVIDER || '';
export const WEBHOOK_PAYLOAD = __ENV.WEBHOOK_PAYLOAD || '';
export const WEBHOOK_HEADERS = __ENV.WEBHOOK_HEADERS || '';

export const AUTH_COMPARE = String(__ENV.AUTH_COMPARE || '').toLowerCase() === 'true';
