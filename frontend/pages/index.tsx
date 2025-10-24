import ShowAllUser from "./components/ShowAllUser";
import { useAuth } from "../context/AuthContext";
import LoginPage from "./components/LoginPage";
import UserDashboard from "./components/UserDashboard";

export default function Home() {
  const { isAuthenticated, role } = useAuth(); // ดึงสถานะมาจาก Context
// 1. ถ้ายังไม่ Login, แสดงหน้า Login
  if (!isAuthenticated) {
    return <LoginPage />;
  }

  // 2. ถ้า Login แล้ว และ Role เป็น 'admin', แสดงหน้า ShowAllUser
  if (role === 'admin') {
    return <ShowAllUser />;
  }

  // 3. ถ้า Login แล้ว แต่ Role ไม่ใช่ 'admin', แสดงหน้า Dashboard ทั่วไป
  return <UserDashboard />;
}