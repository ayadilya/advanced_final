import axios from 'axios';

export const api = axios.create({
    baseURL: 'http://localhost:8080/api', // Replace with your backend API base URL
});

export const registerUser = async (name: string, email: string, password: string) => {
    const response = await api.post('/users', { name, email, password });
    return response.data;
};

export const loginUser = async (email: string, password: string) => {
    const response = await api.post('/users/login', { email, password });
    return response.data;
};

export const fetchProducts = async (token: string) => {
    const response = await api.get('/products', {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
    return response.data; // Ensure this returns an array
};

export const fetchCategories = async (token: string) => {
    const response = await api.get('/categories', {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
    return response.data; // Ensure this returns an array
};

export const getUserInfo = async (token: string) => {
    const response = await api.get('/user/info', {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
    return response.data;
};
