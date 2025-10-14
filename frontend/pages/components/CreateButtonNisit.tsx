import React, { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';

// สร้าง Type สำหรับข้อมูลนิสิตเพื่อให้โค้ดปลอดภัยมากขึ้น
type Nisit = {
  studentCode: string;
  name: string;
  major: string;
};

// ฟังก์ชันสำหรับส่งข้อมูลไปยัง API
const createNisit = async (newNisit: Nisit) => {
  // API ของคุณต้องการข้อมูลในรูปแบบ Array
  const response = await fetch('http://localhost:5000/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    // ส่งข้อมูลนิสิตคนเดียวในรูปแบบ Array
    body: JSON.stringify([newNisit]),
  });

  if (!response.ok) {
    const errorData = await response.json();
    // ดึงข้อความ error จาก API ที่คุณออกแบบไว้
    throw new Error(errorData.error || 'Failed to create nisit');
  }

  return response.json();
};

function CreateButtonNisit() {
  // State สำหรับเปิด/ปิด Modal
  const [isModalOpen, setIsModalOpen] = useState(false);

  // State สำหรับเก็บข้อมูลในฟอร์ม
  const [formData, setFormData] = useState<Nisit>({
    studentCode: '',
    name: '',
    major: '',
  });

  // เข้าถึง QueryClient เพื่อสั่ง refetch ข้อมูลหลังจากการเพิ่มสำเร็จ
  const queryClient = useQueryClient();

  // -----  หัวใจหลัก: การใช้ useMutation  -----
  const mutation = useMutation({
    mutationFn: createNisit,
    onSuccess: () => {
      // เมื่อสำเร็จ:
      alert('เพิ่มข้อมูลนิสิตสำเร็จ!');
      // สั่งให้ queryKey "users" โหลดข้อมูลใหม่ เพื่อให้หน้าเว็บอัปเดต
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setIsModalOpen(false); // ปิด Modal
      setFormData({ studentCode: '', name: '', major: '' }); // ล้างฟอร์ม
    },
    onError: (error) => {
      // เมื่อเกิดข้อผิดพลาด:
      alert(`เกิดข้อผิดพลาด: ${error.message}`);
    },
  });

  // Handler เมื่อมีการเปลี่ยนแปลงใน input
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  // Handler เมื่อกด Submit ฟอร์ม
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.studentCode || !formData.name || !formData.major) {
      alert('กรุณากรอกข้อมูลให้ครบทุกช่อง');
      return;
    }
    // สั่งให้ mutation ทำงานโดยส่งข้อมูลจากฟอร์มไป
    mutation.mutate(formData);
  };

  return (
    <>
      {/* --- ปุ่มสำหรับเปิด Modal --- */}
      <button
        onClick={() => setIsModalOpen(true)}
        className="px-4 py-2 rounded bg-green-600 text-white font-semibold hover:bg-green-700 transition-colors"
      >
        + เพิ่มนิสิต
      </button>

      {/* --- Modal และฟอร์ม (จะแสดงเมื่อ isModalOpen เป็น true) --- */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
          <div className="bg-white p-8 rounded-lg shadow-xl w-full max-w-md">
            <h2 className="text-2xl font-bold mb-6">กรอกข้อมูลนิสิตใหม่</h2>
            
            {/* แสดงข้อความ Error จาก Mutation */}
            {mutation.isError && (
              <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
                <p>{mutation.error.message}</p>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="studentCode" className="block text-sm font-medium text-gray-700">รหัสนักศึกษา</label>
                <input
                  type="text"
                  id="studentCode"
                  name="studentCode"
                  value={formData.studentCode}
                  onChange={handleInputChange}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  required
                />
              </div>
              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-700">ชื่อ-สกุล</label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  required
                />
              </div>
              <div>
                <label htmlFor="major" className="block text-sm font-medium text-gray-700">สาขาวิชา</label>
                <input
                  type="text"
                  id="major"
                  name="major"
                  value={formData.major}
                  onChange={handleInputChange}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  required
                />
              </div>

              {/* --- ปุ่มในฟอร์ม --- */}
              <div className="flex justify-end space-x-4 pt-4">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 rounded bg-gray-200 text-gray-800 hover:bg-gray-300"
                  disabled={mutation.isPending} // Disable ปุ่มขณะกำลังโหลด
                >
                  ยกเลิก
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700 disabled:bg-blue-300"
                  disabled={mutation.isPending} // Disable ปุ่มขณะกำลังโหลด
                >
                  {mutation.isPending ? 'กำลังบันทึก...' : 'บันทึก'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}

export default CreateButtonNisit;