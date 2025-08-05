import { useApp } from "./AppContext"

function GameView() {
    const { gameUrl } = useApp();
    return (
        <div className="w-full h-full flex items-center justify-center py-10 pl-10">
            <div className="w-full h-full">
                {gameUrl.Valid
                    ? <div className="embedded-container w-full h-full">
                        <iframe
                            className="rounded-lg w-full h-full"
                            src={gameUrl.String}
                            title="Embedded Website"
                        />
                    </div>
                    : <div className="embedded-container w-full h-full">
                        <div className="rounded-lg w-full h-full bg-figma-stone flex items-center justify-center">
                            <span className="text-figma-white texg-3xl font-medium">
                                Valid url not provided
                            </span>
                        </div>
                    </div>
                }
            </div>
        </div>
    )
}

export default GameView
