import React from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrash } from '@fortawesome/free-solid-svg-icons';
type Nisit = {
    studentCode: string;
    name: string;
    major: string;
};

type DeleteButtonNisitProps = {
    initialData: Nisit;
};

// 1. แก้ไขฟังก์ชัน API call ให้รับเฉพาะ studentCode
// และส่งไปเป็น Query Parameter
const deleteNisit = async (studentCode: string) => {
    const response = await fetch(`http://localhost:5000/users?studentCode=${studentCode}`, {
        method: 'DELETE',
    });

    if (!response.ok) {
        // ลองอ่าน error message จาก backend
        const errorData = await response.json().catch(() => ({ error: 'An unknown error occurred' }));
        throw new Error(errorData.error || 'Failed to delete nisit');
    }
    return response.json();
};

function DeleteButtonNisit({ initialData }: DeleteButtonNisitProps) {
    const queryClient = useQueryClient();

    const mutation = useMutation({
        // 2. ส่งฟังก์ชันที่แก้ไขแล้วเข้ามา
        mutationFn: deleteNisit,
        onSuccess: () => {
            alert('Delete nisit succesful!');
            queryClient.invalidateQueries({ queryKey: ['users'] });
        },
        onError: (error) => {
            alert(`Can't delete nisit: ${error.message}`);
        },
    });

    // 3. สร้างฟังก์ชัน handleDelete ที่ชัดเจน
    const handleDelete = () => {
        if (window.confirm(`Are you sure to delete ${initialData.name}?`)) {
            mutation.mutate(initialData.studentCode);
        }
    };

    return (
        <button
            onClick={handleDelete}
            className="px-3 py-1 text-sm rounded bg-gray-300 text-white font-semibold hover:bg-gray-500 transition-colors disabled:bg-red-300"
            disabled={mutation.isPending}
        >
            <FontAwesomeIcon icon={faTrash} />
        </button>
    );
}

export default DeleteButtonNisit;