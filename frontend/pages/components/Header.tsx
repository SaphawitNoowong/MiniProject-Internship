import React, { useState } from 'react';
import { FaBars, FaTimes, FaUserCircle } from 'react-icons/fa'; // ไอคอนสำหรับเมนู

function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const navLinks = [
    { title: 'หน้าหลัก', href: '/dashboard' },
    { title: 'ตารางเรียน', href: '/schedule' },
    { title: 'ผลการเรียน', href: '/grades' },
    { title: 'ลงทะเบียนเรียน', href: '/courses' },
  ];

  return (
    <header className="bg-white shadow-md sticky top-0 z-50">
      {/* Main Header Bar */}
      <div className="container mx-auto flex justify-between items-center p-4">
        {/* Logo and Site Title */}
        <div className="flex items-center space-x-3">
          <img src="/public/images/WUlogo.jpg" alt="" className="h-10" />
          <a href="/" className="text-xl font-bold text-gray-800 hover:text-blue-600">
            WU University
          </a>
        </div>
      </div>
    </header>
  );
}

export default Header;