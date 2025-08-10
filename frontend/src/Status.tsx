import { useApp } from "./AppContext";


const Status = () => {
    const {
        loading,
        notSaved,
    } = useApp();

    if (loading) return (
        <div className="flex gap-1 items-center">
            <div className="rounded-full size-2 bg-orange-600"></div>
            <span className="text-orange-600 text-sm font-medium">
                Pending...
            </span>
        </div>
    )

    if (notSaved) return (
        <div className="flex gap-1 items-center">
            <div className="rounded-full size-2 bg-red-600"></div>
            <span className="text-red-600 text-sm font-medium">
                Not Saved
            </span>
        </div>
    )

    return (
        <div className="flex gap-1 items-center">
            <div className="rounded-full size-2 bg-green-700"></div>
            <span className="text-green-700 text-sm font-medium">
                Saved
            </span>
        </div>
    );
}

export default Status;
