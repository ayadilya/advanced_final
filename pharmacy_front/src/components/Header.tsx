"use client";
import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { getUserInfo } from '@/services/api'; // Correctly importing getUserInfo

const Header: React.FC = () => {
    const [user, setUser] = useState<string | null>(null);
    const router = useRouter();

    useEffect(() => {
        const fetchUser = async () => {
            const token = localStorage.getItem('token');
            if (token) {
                try {
                    const userInfo = await getUserInfo(token);
                    setUser(userInfo.name);
                } catch (error) {
                    console.error('Failed to fetch user info', error);
                }
            }
        };

        fetchUser();
    }, [router]);

    const handleLogout = () => {
        localStorage.removeItem('token');
        setUser(null);
        router.push('/');
    };

    return (
        <header className="bg-blue-500 text-white p-4">
            <div className="container mx-auto flex justify-between items-center">
                <h1 className="text-2xl font-bold">Pharmacy Store</h1>
                <nav>
                    {user ? (
                        <div className="flex items-center">
                            <span className="mr-4">Hello, {user}</span>
                            <Link href="/create-category">
                                <span className="mr-4">Create Category</span>
                            </Link>
                            <Link href="/create-product">
                                <span className="mr-4">Create Product</span>
                            </Link>
                            <button
                                onClick={handleLogout}
                                className="bg-red-500 text-white px-4 py-2 rounded"
                            >
                                Logout
                            </button>
                        </div>
                    ) : (
                        <div className="flex items-center">
                            <Link href="/login">
                                <span className="mr-4">Login</span>
                            </Link>
                            <Link href="/register">
                                <span>Register</span>
                            </Link>
                        </div>
                    )}
                </nav>
            </div>
        </header>
    );
};

export default Header;
