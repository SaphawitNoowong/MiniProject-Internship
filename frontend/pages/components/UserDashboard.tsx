import React from 'react'
import { useAuth } from '../../context/AuthContext';
import Header from './Header';
function UserDashboard() {
    const { logout } = useAuth();
    return (
        <div className="bg-gray-50 min-h-screen">
            <Header />
            <h1 className="text-2xl">Welcome, User!</h1>
            <p>This is your dashboard.</p>
        </div>
    );
}

export default UserDashboard