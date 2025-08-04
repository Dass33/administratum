import { useApp } from "./AppContext"

function GameView() {
    const { gameUrl } = useApp();
    return (
        <div className="w-full h-full flex items-center justify-center py-10 pl-10">
            <div className="w-full h-full">
                <div className="embedded-container w-full h-full">
                    <iframe
                        className="rounded-lg w-full h-full"
                        src={gameUrl}
                        title="Embedded Website"
                    />
                </div>
            </div>
        </div>
    )
}

export default GameView
