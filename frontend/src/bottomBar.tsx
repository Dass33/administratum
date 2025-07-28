import { useEffect, useState } from "react";
import settings from "./assets/settings.svg"
import Dropdown from "./dropdown";
import { useApp, CurrSheet, ColSuffix } from "./AppContext";

const BottomBar = () => {
    const {
        sheets, setSheets,
        currSheet, setCurrSheet,
        setCurrTable,
        columns, setColumns,
    } = useApp();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
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
        <div className="flex flex-row gap-5 items-center">
            <button className="hover:scale-110 transition-transform duration-100">
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
        </div>
    );
}

export default BottomBar;

