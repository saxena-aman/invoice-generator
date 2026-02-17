const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

/**
 * Make an auth API request and return parsed JSON.
 * Throws an Error with the server message on failure.
 */
async function authFetch(endpoint, body) {
    const res = await fetch(`${API_URL}/auth${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
    });

    const data = await res.json();

    if (!res.ok) {
        throw new Error(data.message || 'Something went wrong');
    }

    return data; // TokenResponse { accessToken, refreshToken, expiresIn, tokenType }
}

export function loginUser(email, password) {
    return authFetch('/login', { email, password });
}

export function registerUser(email, password, name) {
    return authFetch('/register', { email, password, name });
}

export function refreshToken(refreshToken) {
    return authFetch('/refresh', { refreshToken });
}
