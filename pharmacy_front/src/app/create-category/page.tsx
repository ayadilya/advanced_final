"use client";
import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/services/api'; // Correctly importing api

const CreateCategory: React.FC = () => {
    const [name, setName] = useState('');
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const token = localStorage.getItem('token');
            if (token) {
                await api.post('/categories', { name }, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                router.push('/');
            }
        } catch (error) {
            console.error('Failed to create category', error);
        }
    };

    return (
        <div className="container mx-auto p-4">
            <h2 className="text-3xl font-bold mb-4">Create Category</h2>
            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label className="block text-gray-700">Name</label>
                    <input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="w-full p-2 border border-gray-300 rounded"
                        required
                    />
                </div>
                <button type="submit" className="bg-blue-500 text-white p-2 rounded">
                    Create
                </button>
            </form>
        </div>
    );
};

export default CreateCategory;
