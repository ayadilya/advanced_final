"use client";
import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { fetchCategories, api } from '@/services/api'; // Correctly importing api and fetchCategories

const CreateProduct: React.FC = () => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [price, setPrice] = useState('');
    const [stock, setStock] = useState('');
    const [categoryID, setCategoryID] = useState('');
    const [categories, setCategories] = useState<any[]>([]); // Ensure the type is an array
    const router = useRouter();

    useEffect(() => {
        const loadCategories = async () => {
            const token = localStorage.getItem('token');
            if (token) {
                try {
                    const data = await fetchCategories(token);
                    console.log('Fetched categories', data);
                    setCategories(Array.isArray(data.message) ? data.message : []); // Ensure data is an array
                } catch (error) {
                    console.error('Failed to fetch categories', error);
                }
            }
        };

        loadCategories();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const token = localStorage.getItem('token');
            if (token) {
                await api.post('/products', {
                    name,
                    description,
                    price: parseFloat(price),
                    stock: parseInt(stock),
                    category_id: parseInt(categoryID),
                }, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                router.push('/');
            }
        } catch (error) {
            console.error('Failed to create product', error);
        }
    };

    return (
        <div className="container mx-auto p-4">
            <h2 className="text-3xl font-bold mb-4">Create Product</h2>
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
                <div className="mb-4">
                    <label className="block text-gray-700">Description</label>
                    <textarea
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        className="w-full p-2 border border-gray-300 rounded"
                        required
                    ></textarea>
                </div>
                <div className="mb-4">
                    <label className="block text-gray-700">Price</label>
                    <input
                        type="number"
                        step="0.01"
                        value={price}
                        onChange={(e) => setPrice(e.target.value)}
                        className="w-full p-2 border border-gray-300 rounded"
                        required
                    />
                </div>
                <div className="mb-4">
                    <label className="block text-gray-700">Stock</label>
                    <input
                        type="number"
                        value={stock}
                        onChange={(e) => setStock(e.target.value)}
                        className="w-full p-2 border border-gray-300 rounded"
                        required
                    />
                </div>
                <div className="mb-4">
                    <label className="block text-gray-700">Category</label>
                    <select
                        value={categoryID}
                        onChange={(e) => setCategoryID(e.target.value)}
                        className="w-full p-2 border border-gray-300 rounded"
                        required
                    >
                        <option value="">Select a category</option>
                        {categories.map((category) => (
                            <option key={category.id} value={category.id}>
                                {category.name}
                            </option>
                        ))}
                    </select>
                </div>
                <button type="submit" className="bg-blue-500 text-white p-2 rounded">
                    Create
                </button>
            </form>
        </div>
    );
};

export default CreateProduct;
