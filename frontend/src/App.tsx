import BottomBar from "./bottomBar";
import GameView from "./gameView"
import SaveButton from "./saveButton";
import SelectBranch from "./selectBranch";
import SelectProject from "./selectProject";
import Table from "./table";

GameView
function App() {

    return (
        <div className='bg-figma-white h-screen pt-10 px-10 flex flex-row'>
            <div className="flex flex-col flex-grow">
                <div className="flex flex-row">
                    <SaveButton />

                    <div className="flex items-center text-lg">
                        <SelectProject />
                        <span className="mx-2">/</span>
                        <SelectBranch />
                    </div>
                </div>
                <div className="mt-10">
                    <Table />
                </div>
                <div className="bottom-0 absolute my-2">
                    <BottomBar />
                </div>
            </div>
            <GameView />
        </div>
    )
}

export default App
