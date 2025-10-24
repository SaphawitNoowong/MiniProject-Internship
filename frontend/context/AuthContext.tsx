import React, { createContext, useContext, useState, useEffect } from 'react';

// 1. สร้าง Type สำหรับ State
interface AuthContextType {
  isAuthenticated: boolean;
  role: string | null;
  login: (studentCode: string, password: string) => Promise<void>;
  logout: () => void;
}

// 2. สร้าง Context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// 3. สร้าง Provider (ตัวจัดการ State)
export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [role, setRole] = useState<string | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

  // 4. (สำคัญ) เมื่อ App โหลด, ให้เช็ค localStorage ว่าเคย Login ค้างไว้ไหม
  useEffect(() => {
    const storedRole = localStorage.getItem('userRole');
    if (storedRole) {
      setRole(storedRole);
      setIsAuthenticated(true);
    }
  }, []);

  // 5. ฟังก์ชัน Login
  const login = async (studentCode: string, password: string) => {
    const response = await fetch('http://localhost:5000/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ studentCode, password }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Login failed');
    }

    const data = await response.json();
    
    // 6. บันทึก Role ลง State และ localStorage
    setRole(data.role);
    setIsAuthenticated(true);
    localStorage.setItem('userRole', data.role);
  };

  // 7. ฟังก์ชัน Logout
  const logout = () => {
    setRole(null);
    setIsAuthenticated(false);
    localStorage.removeItem('userRole');
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, role, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// 8. สร้าง Hook สำหรับเรียกใช้ง่ายๆ
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};