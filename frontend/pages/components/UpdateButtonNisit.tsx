import React, { useState, useEffect } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';

type Nisit = {
    studentCode: string;
    name: string;
    major: string;
};

type UpdateButtonNisitProps = {
    initialData: Nisit;
};

const updateNisit = async (updatedNisit: Nisit) => {
    const response = await fetch('http://localhost:5000/users', {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify([updatedNisit]),
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to update nisit');
    }
    return response.json();
};


// 2. รับ props เข้ามาใน Component
function UpdateButtonNisit({ initialData }: UpdateButtonNisitProps) {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [formData, setFormData] = useState<Nisit>(initialData);

    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: updateNisit,
        onSuccess: () => {
            alert('Update nisit successful!');
            queryClient.invalidateQueries({ queryKey: ['users'] });
            setIsModalOpen(false);
        },
        onError: (error) => {
            alert(`Can't update nisit :${error.message}`);
        },
    });

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const payload = {
            studentCode: formData.studentCode,
            name: formData.name,
            major: formData.major,
        };
        mutation.mutate(payload);
};

// ฟังก์ชันสำหรับเปิด Modal และตั้งค่า form ใหม่ทุกครั้ง
// เพื่อป้องกันกรณีที่ props เปลี่ยน แต่ form ยังเป็นข้อมูลเก่า
const openModal = () => {
    setFormData(initialData);
    setIsModalOpen(true);
};

return (
    <>
        <button
            onClick={openModal}
            className="px-3 py-1 text-sm rounded bg-yellow-500 text-white font-semibold hover:bg-yellow-600 transition-colors"
        >
            Edit
        </button>

        {isModalOpen && (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
                <div className="bg-white p-8 rounded-lg shadow-xl w-full max-w-md">
                    <h2 className="text-2xl font-bold mb-6">Update nisit</h2>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <label htmlFor="studentCode" className="block text-sm font-medium text-gray-700">Student Code (Can't edit)</label>
                            <input
                                type="text"
                                id="studentCode"
                                name="studentCode"
                                value={formData.studentCode}
                                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm bg-gray-100 text-gray-500"
                                readOnly
                            />
                        </div>
                        <div>
                            <label htmlFor="name" className="block text-sm font-medium text-gray-700">Name-lastname</label>
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
                            <label htmlFor="major" className="block text-sm font-medium text-gray-700">Major</label>
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
                        <div className="flex justify-end space-x-4 pt-4">
                            <button
                                type="button"
                                onClick={() => setIsModalOpen(false)}
                                className="px-4 py-2 rounded bg-gray-200 text-gray-800 hover:bg-gray-300"
                                disabled={mutation.isPending}
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                className="px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700 disabled:bg-blue-300"
                                disabled={mutation.isPending}
                            >
                                {mutation.isPending ? 'Confirming...' : 'Confirm'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        )}
    </>
);
}

export default UpdateButtonNisit;