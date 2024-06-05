"use client";
import type { NextPage } from 'next';
import Head from 'next/head';
import { useEffect, useState } from 'react';
import Header from '../components/Header';
import Footer from '../components/Footer';
import { fetchProducts } from '../services/api';

const Home: NextPage = () => {
  const [products, setProducts] = useState<any[]>([]); // Ensure the type is an array

  useEffect(() => {
    const loadProducts = async () => {
      const token = localStorage.getItem('token');
      if (token) {
        try {
          const data = await fetchProducts(token);
          setProducts(Array.isArray(data.message) ? data.message : []); // Ensure data is an array
          console.log('Fetched products', data);
        } catch (error) {
          console.error('Failed to fetch products', error);
        }
      }
    };

    loadProducts();
  }, []);

  return (
    <div className="flex flex-col min-h-screen">
      <Head>
        <title>Pharmacy Store</title>
        <meta name="description" content="Welcome to the Pharmacy Store" />
      </Head>

      <Header />

      <main className="flex-grow container mx-auto p-4">
        <h2 className="text-3xl font-bold mb-4">Welcome to the Pharmacy Store</h2>
        <p className="mb-4">Browse our products and find what you need for your health and wellness.</p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {products.map((product) => (
            <div key={product.id} className="border p-4 rounded">
              <h3 className="text-xl font-bold mb-2">{product.name}</h3>
              <p className="mb-2">{product.description}</p>
              <p className="font-bold">${product.price}</p>
            </div>
          ))}
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default Home;
