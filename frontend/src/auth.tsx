import React, { useState } from 'react';
import { useApp, LoginData, isValidEmail, Domain } from './AppContext';
import logo from "./assets/logo.svg";
import reload_black from "./assets/reload_black.svg";
import reload from "./assets/reload.svg";

type LoginCredentials = {
    email: string;
    password: string;
};

enum AuthAction {
    Login = "/login",
    Register = "/register"
}

const Auth = () => {
    const {
        setAuthenticated,
        setAccessToken,
        setColumns,
        setCurrSheet,
        setCurrTable,
        setTableNames,
        setGameUrl,
        setCurrBranch,
    } = useApp();
    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const [email, setEmail] = useState<string>();
    const [password, setPassword] = useState<string>();
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<string | null>(null);

    const handleSubmit = async (data: LoginCredentials, type: string) => {
        if (!isValidEmail(data.email)) {
            setError('Please enter a valid email address');
            return;
        }

        setLoading(type);
        try {
            const response = await fetch(Domain + type, {
                method: 'POST',
                credentials: "include",
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                setError(`HTTP error! status: ${response.status}`);
                setLoading(null);
                return;
            }
            const result: LoginData = await response.json();
            console.log(result);
            setCurrSheet(result.opened_sheet)
            setCurrTable(result.opened_table)
            setColumns(result.opened_sheet.columns);
            setTableNames(result.table_names)
            setError(null);
            setAccessToken(result.token)
            setAuthenticated(true);
            setGameUrl(result.opened_table.game_url)
            setCurrBranch(result.opened_sheet.curr_branch);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            setLoading(null);
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
                    {error && <span className='text-red-600 text-xs ml-1'>
                        {error === 'Incorrect credentials' ? 'Incorrect credentials' : error}
                    </span>}


                    <div className={`flex gap-4 justify-between ${error ? "mt-2" : "mt-5"}`}>
                        <button className='bg-figma-green w-20 rounded-lg p-2 px-4 text-figma-white font-bold'
                            onClick={() => {
                                if (email && password) {
                                    handleSubmit({ email: email, password: password }, AuthAction.Login)
                                }
                            }}
                        >
                            {loading == AuthAction.Login
                                ? <img className="size-6 animate-spin mx-auto" src={reload} />
                                : <span>Login</span>
                            }
                        </button>
                        <button className='rounded-lg p-2 px-4 w-20 text-figma-black'
                            onClick={() => {
                                if (email && password) {
                                    handleSubmit({ email: email, password: password }, AuthAction.Register)
                                }
                            }}
                        >
                            {loading == AuthAction.Register
                                ? <img className="size-6 mx-auto mr-1 animate-spin" src={reload_black} />
                                : <span>Register</span>
                            }
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Auth;
