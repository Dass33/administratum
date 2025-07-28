import { useEffect, useState } from "react";
import settings from "./assets/settings.svg"
import Dropdown from "./dropdown";
import { useApp, CurrSheet, ColSuffix, Sheets } from "./AppContext";
import plus from "./assets/plus.svg"

const BottomBar = () => {
    const {
        sheets, setSheets,
        currSheet, setCurrSheet,
        setCurrTable,
        columns, setColumns,
        setSheetModal,
        setSettingsModal,
    } = useApp();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const stored = localStorage.getItem(Sheets);
        console.log(stored);
        if (stored) {
            setSheets(JSON.parse(stored))
            setLoading(false);
            return
        }

        fetch('http://localhost:8080/sheets')
            .then(response => response.json())
            .then((data: string[]) => {
                setSheets(data);
                if (!currSheet) setCurrSheet(data[0])
                setLoading(false);
            })
            .catch(error => {
                setError(error);
                setLoading(false);
            });
    }, []);

    return (
        <div className="flex flex-row gap-4 items-center">
            <button className="hover:scale-110 transition-transform duration-100 mr-1"
                onClick={() => setSettingsModal(true)}
            >
                <img className="" src={settings} />
            </button>
            {(!loading && !error) &&
                <Dropdown
                    options={sheets.map(item => {
                        return { value: item, label: item }
                    })}
                    placeholder={currSheet}
                    onSelect={(item) => {
                        localStorage.setItem(currSheet + ColSuffix, JSON.stringify(columns));
                        setCurrSheet(item.value)
                        localStorage.setItem(CurrSheet, item.value);
                        setColumns(() => {
                            const stored = localStorage.getItem(item.value + ColSuffix);
                            return stored ? JSON.parse(stored) : [];
                        })
                        setCurrTable(() => {
                            const stored = localStorage.getItem(item.value);
                            return stored ? JSON.parse(stored) : [];
                        });
                    }}
                    isDown={false}
                />
            }
            <button className="hover:scale-125 transition-transform duration-100"
                onClick={() => setSheetModal(true)}>
                <img src={plus} />
            </button>
        </div>
    );
}

export default BottomBar;

