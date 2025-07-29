import React, { useState } from 'react';
import { useApp } from './AppContext';
import logo from "./assets/logo.svg";

type loginData = {
    email: string;
    password: string;
};

const Auth = () => {
    const {
        setAuthenticated
    } = useApp();
    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const [email, setEmail] = useState<string>();
    const [password, setPassword] = useState<string>();
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (data: loginData, type: string) => {
        try {
            const response = await fetch('http://localhost:8080/' + type, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                setError(`HTTP error! status: ${response.status}`);
                return;
            }

            const result = await response.json();
            setError(null);
            console.log('Success:', result);
            setAuthenticated(true);

        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
        }
    };

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
        >
            <div onClick={stopPropagation}>
                <div className="p-4 rounded-lg text-figma-black bg-figma-white font-medium
                    overflow-y-auto focus:outline-none w-[22rem] mb-16">
                    <div className='flex items-center mb-4 gap-4'>
                        <img className='size-12 border-2 bg-figma-white border-figma-black p-1 rounded-lg' src={logo} />
                        <h2 className='text-2xl text-figma-black'>Login</h2>
                    </div>
                    <div className=''>
                        <input className="w-80 border border-figma-gray bg-figma-white focus:outline-none rounded-lg p-2"
                            type='email'
                            placeholder="Email"
                            onChange={(e) => { setEmail(e.target.value) }}
                        />
                    </div>

                    <div className='mt-4'>
                        <input className="w-80 border border-figma-gray bg-figma-white focus:outline-none rounded-lg p-2"
                            type='password'
                            placeholder="Password"
                            onChange={(e) => { setPassword(e.target.value) }}
                        />
                    </div>
                    {/*<button className='rounded-lg mt-2 text-xs underline ml-1 text-figma-black'>
                        Reset password
                    </button>*/}
                    {error && <span className='text-red-600 text-xs ml-1'>Incorrect credentials</span>}


                    <div className={`flex gap-4 justify-between ${error ? "mt-2" : "mt-5"}`}>
                        <button className='bg-figma-green rounded-lg p-2 px-4 text-figma-white font-bold'
                            onClick={() => {
                                if (email && password) {
                                    handleSubmit({ email: email, password: password }, "login")
                                }
                            }}
                        >
                            Login
                        </button>
                        <button className='rounded-lg p-2 px-4 text-figma-black'
                            onClick={() => {
                                if (email && password) {
                                    handleSubmit({ email: email, password: password }, "register")
                                }
                            }}
                        >
                            Register
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Auth;
