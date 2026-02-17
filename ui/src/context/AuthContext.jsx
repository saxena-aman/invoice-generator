import { createContext, useContext, useState, useCallback, useMemo } from 'react';
import { loginUser, registerUser } from '../utils/authApi';

const AuthContext = createContext(null);

function loadPersistedAuth() {
  try {
    const token = localStorage.getItem('auth_token');
    const user = localStorage.getItem('auth_user');
    if (token && user) {
      return { token, user: JSON.parse(user) };
    }
  } catch { /* ignore corrupt data */ }
  return { token: null, user: null };
}

export const AuthProvider = ({ children }) => {
  const [authState, setAuthState] = useState(loadPersistedAuth);

  const persistAuth = useCallback((token, user) => {
    localStorage.setItem('auth_token', token);
    localStorage.setItem('auth_user', JSON.stringify(user));
    setAuthState({ token, user });
  }, []);

  const login = useCallback(async (email, password) => {
    const data = await loginUser(email, password);
    // We don't get user info from login response directly, so decode from token or store email
    const user = { email };
    persistAuth(data.accessToken, user);
    // Store refresh token separately
    localStorage.setItem('auth_refresh_token', data.refreshToken);
    return data;
  }, [persistAuth]);

  const signup = useCallback(async (email, password, name) => {
    const data = await registerUser(email, password, name);
    const user = { email, name };
    persistAuth(data.accessToken, user);
    localStorage.setItem('auth_refresh_token', data.refreshToken);
    return data;
  }, [persistAuth]);

  const logout = useCallback(() => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('auth_user');
    localStorage.removeItem('auth_refresh_token');
    setAuthState({ token: null, user: null });
  }, []);

  const value = useMemo(() => ({
    user: authState.user,
    token: authState.token,
    isAuthenticated: !!authState.token,
    login,
    signup,
    logout,
  }), [authState, login, signup, logout]);

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
