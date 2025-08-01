import BottomBar from "./bottomBar";
import GameView from "./gameView"
import SaveButton from "./saveButton";
import SelectBranch from "./selectBranch";
import SelectProject from "./selectProject";
import Permissions from "./permissions";
import Table from "./table";
import { useApp } from "./AppContext";
import CellModal from "./cellModal";
import ColModal from "./colModal";
import SheetModal from "./sheetModal";
import SettingsModal from "./settingsModal";
import Auth from "./auth";

function App() {
    const {
        cellModal,
        colModal,
        sheetModal,
        settingsModal,
        authenticated,
        loading,
    } = useApp()

    if (loading) return (
        <div className='bg-figma-white h-screen pl-10 xl:pr-5 flex items-center'>
            <h1 className="text-figma-black text-2xl text-center w-full">Loading...</h1>
        </div>
    )

    if (!authenticated) return (
        <div className='bg-figma-white h-screen pl-10 xl:pr-5 flex flex-row'>
            <Auth />
        </div>
    )

    return (
        <div className='bg-figma-white h-screen pl-10 xl:pr-5 flex flex-row'>
            {cellModal && <CellModal />}
            {colModal > -1 && <ColModal />}
            {sheetModal && <SheetModal />}
            {settingsModal && <SettingsModal />}
            <div className="flex flex-col flex-grow my-10">
                <div className="flex flex-row gap-4">
                    <SaveButton />

                    <div className="flex items-center text-lg">
                        <SelectProject />
                        <span className="mx-2">/</span>
                        <SelectBranch />
                    </div>
                    <Permissions />
                </div>
                <div className="mt-10 grow">
                    <Table />
                </div>
                <div className="bottom-1 absolute my-2">
                    <BottomBar />
                </div>
            </div>
            <GameView />
        </div>
    )
}

export default App
