import React, { useEffect, useState } from 'react';
import { useApp, PermissionsEnum, isValidEmail, Domain } from './AppContext';
import Dropdown from './dropdown';

const ShareModal = () => {
    const {
        setShareModal,
        accessToken,
        currTable,
    } = useApp();

    const [email, setEmail] = useState<string>();
    const [selectedPerm, setSelectedPerm] = useState<string>(PermissionsEnum.CONTRIBUTOR);

    const permissionsOptions = Object.values(PermissionsEnum).map(item => {
        return { label: item, value: item }
    })
    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const saveAndExit = () => {
        setShareModal(false);
    }

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') saveAndExit()
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [setShareModal]);

    const handleShare = () => {
        if (!email || !accessToken || !currTable || !isValidEmail(email)) return;
        postNewShare(email, selectedPerm, currTable.id, accessToken);
        setShareModal(false);
    };

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div onClick={stopPropagation}>
                <div className="p-4 rounded-lg text-figma-black bg-figma-white font-medium
                    overflow-y-auto focus:outline-none w-[35rem]">
                    <h2 className='text-2xl mb-6'>Share Project</h2>

                    <div className='flex gap-4 mb-6 w-full'>
                        <input className={`grow border bg-figma-white focus:outline-none rounded-lg p-2
                                        ${isValidEmail(email) || !email?.length ? "border-figma-gray" : "border-red-600"}`}
                            type='email'
                            placeholder="Email"
                            onChange={(e) => { setEmail(e.target.value) }}
                        />

                        <Dropdown
                            options={permissionsOptions}
                            defaultValue={PermissionsEnum.CONTRIBUTOR}
                            onSelect={(val) => setSelectedPerm(val.value)}
                        />
                    </div>

                    <div className='w-full flex justify-end'>
                        <button className='bg-green-600 w-24 rounded-lg p-2 px-4 text-figma-white font-bold mt-4 hover:bg-green-700'
                            onClick={handleShare}
                        >
                            <span>Share</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const postNewShare = (email: string, perm: string, tableId: string, token: string) => {
    const newShareParams: { email: string, perm: string, table_id: string } = {
        email: email,
        perm: perm,
        table_id: tableId,
    };

    fetch(Domain + '/add_share', {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(newShareParams)
    })
        .then(response => {
            if (response.status != 201) {
                throw "Could not share"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

export default ShareModal; 
