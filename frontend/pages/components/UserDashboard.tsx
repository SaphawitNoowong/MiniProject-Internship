import React from 'react'
import { useAuth } from '../../context/AuthContext';

function UserDashboard() {
    const { logout } = useAuth();
    return (
        <div className="p-8">
            <h1 className="text-2xl">Welcome, User!</h1>
            <p>This is your dashboard.</p>
            <button onClick={logout} className="mt-4 px-4 py-2 bg-red-500 text-white rounded">
                Logout
            </button>
        </div>
    );
}

export default UserDashboard