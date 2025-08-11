import { useState, useEffect } from "react";
import { useApp, Domain, Sheet } from "./AppContext";
import Dropdown, { DropdownOption } from './dropdown';

type MergeConflict = {
    id: string;
    type: string;
    sheet_id: string;
    sheet_name: string;
    column_id?: string;
    column_name?: string;
    row_index?: number;
    property?: string;
    source_value: string;
    target_value: string;
    source_updated_at: string;
    target_updated_at: string;
};

type MergePreviewResponse = {
    conflicts: MergeConflict[];
};

type MergeResolution = {
    conflict_id: string;
    chosen_source: string;
};

export default function MergeModal() {
    const {
        currTable,
        currSheet,
        accessToken,
        setMergeModal,
        setCurrSheet,
        setColumns
    } = useApp();

    const [selectedBranch, setSelectedBranch] = useState<string>("");
    const [conflicts, setConflicts] = useState<MergeConflict[]>([]);
    const [resolutions, setResolutions] = useState<Record<string, string>>({});
    const [loading, setLoading] = useState(false);
    const [step, setStep] = useState<'select' | 'conflicts' | 'no-conflicts' | 'merging'>('select');
    const [error, setError] = useState<string>("");

    const availableBranches = currTable?.branches_id_names?.filter(
        branch => branch.id !== currSheet?.curr_branch.id
    ) || [];

    const reloadCurrentSheet = async () => {
        if (!currSheet || !accessToken) return;

        try {
            const response = await fetch(`${Domain}/get_sheet/${currSheet.id}`, {
                method: "GET",
                headers: {
                    'Authorization': `Bearer ${accessToken}`
                },
                credentials: "include"
            });

            if (!response.ok) {
                throw new Error("Could not reload sheet data");
            }

            const result: Sheet = await response.json();
            console.log(result);
            setCurrSheet(result);
            setColumns(result.columns);
        } catch (error) {
            console.error("Failed to reload sheet:", error);
        }
    };

    const handlePreview = async () => {
        if (!selectedBranch || !currSheet) return;

        console.log('=== MERGE PREVIEW START ===');
        console.log('Source branch ID:', selectedBranch);
        console.log('Target branch ID:', currSheet.curr_branch.id);
        console.log('Current sheet:', currSheet);

        setLoading(true);
        setError("");

        try {
            const requestBody = {
                source_branch_id: selectedBranch,
                target_branch_id: currSheet.curr_branch.id
            };
            console.log('Request body:', requestBody);

            const response = await fetch(`${Domain}/merge_preview`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify(requestBody)
            });

            console.log('Response status:', response.status);
            console.log('Response ok:', response.ok);

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Error response:', errorText);
                throw new Error(`Failed to preview merge: ${response.status}`);
            }

            const data: MergePreviewResponse = await response.json();
            console.log('Preview response data:', data);
            const conflicts = data.conflicts || [];
            console.log('Conflicts found:', conflicts.length);
            setConflicts(conflicts);

            if (conflicts.length === 0) {
                console.log('No conflicts - moving to no-conflicts step');
                setStep('no-conflicts');
            } else {
                console.log('Conflicts found - moving to conflicts step');
                setStep('conflicts');
            }
        } catch (error) {
            setError(error instanceof Error ? error.message : 'Unknown error occurred');
        } finally {
            setLoading(false);
        }
    };

    const handleMerge = async () => {
        if (!selectedBranch || !currSheet) return;

        console.log('=== MERGE EXECUTE START ===');
        console.log('Source branch ID:', selectedBranch);
        console.log('Target branch ID:', currSheet.curr_branch.id);
        console.log('Conflicts to resolve:', conflicts.length);
        console.log('Current resolutions:', resolutions);

        setStep('merging');
        setLoading(true);
        setError("");

        const mergeResolutions: MergeResolution[] = conflicts.map(conflict => ({
            conflict_id: conflict.id,
            chosen_source: resolutions[conflict.id] || 'target'
        }));

        console.log('Final merge resolutions:', mergeResolutions);

        try {
            const requestBody = {
                source_branch_id: selectedBranch,
                target_branch_id: currSheet.curr_branch.id,
                resolutions: mergeResolutions
            };
            console.log('Merge request body:', requestBody);

            const response = await fetch(`${Domain}/merge_execute`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify(requestBody)
            });

            console.log('Merge response status:', response.status);
            console.log('Merge response ok:', response.ok);

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Merge error response:', errorText);
                throw new Error(`Failed to execute merge: ${response.status}`);
            }

            const responseData = await response.json();
            console.log('Merge response data:', responseData);

            console.log('Reloading current sheet...');
            await reloadCurrentSheet();
            console.log('Sheet reloaded, closing modal');
            setMergeModal(false);
        } catch (error) {
            setError(error instanceof Error ? error.message : 'Failed to merge branches');
            setStep('conflicts');
        } finally {
            setLoading(false);
        }
    };

    const handleResolution = (conflictId: string, source: string) => {
        setResolutions(prev => ({
            ...prev,
            [conflictId]: source
        }));
    };

    const allConflictsResolved = !conflicts || conflicts.length === 0 || conflicts.every(conflict =>
        resolutions[conflict.id]
    );

    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const saveAndExit = () => {
        setMergeModal(false);
    }

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') saveAndExit()
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [setMergeModal]);

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div onClick={stopPropagation}>
                <div className="p-4 rounded-lg text-figma-black bg-figma-white font-medium overflow-y-auto focus:outline-none w-[35rem]">
                    <h2 className='text-2xl mb-4'>Merge Branch</h2>

                    {error && (
                        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded-lg">
                            {error}
                        </div>
                    )}

                    {step === 'select' && (
                        <div>
                            <p className="mb-4">
                                Select a branch to merge into <strong>{currSheet?.curr_branch.name}</strong>:
                            </p>

                            <div className="flex justify-between mt-8">
                                <Dropdown
                                    options={availableBranches.map(branch => ({ label: branch.name, value: branch.id }))}
                                    placeholder="Select a branch..."
                                    onSelect={(option: DropdownOption) => setSelectedBranch(option.value)}
                                    usePortal={true}
                                />

                                <button
                                    onClick={handlePreview}
                                    disabled={!selectedBranch || loading}
                                    className="px-4 py-2 bg-green-600 text-figma-white font-medium rounded-lg hover:bg-green-700 disabled:opacity-50"
                                >
                                    {loading ? 'Checking...' : 'Preview Merge'}
                                </button>
                            </div>
                        </div>
                    )}

                    {step === 'conflicts' && (
                        <div>
                            <p className="mb-4">
                                Found {conflicts?.length || 0} conflict{(conflicts?.length || 0) !== 1 ? 's' : ''} that need resolution:
                            </p>

                            <div className="space-y-4 max-h-96 overflow-y-auto">
                                {conflicts?.map(conflict => (
                                    <div key={conflict.id} className="border border-figma-gray rounded-lg p-4 bg-figma-white">
                                        <div className="mb-2">
                                            <strong>
                                                {conflict.type === 'cell_data' && `Cell conflict in ${conflict.sheet_name}/${conflict.column_name} (row ${conflict.row_index})`}
                                                {conflict.type === 'column_property' && `Column ${conflict.property} conflict in ${conflict.sheet_name}/${conflict.column_name}`}
                                                {conflict.type === 'sheet_property' && `Sheet ${conflict.property} conflict in ${conflict.sheet_name}`}
                                            </strong>
                                        </div>

                                        <div className="grid grid-cols-2 gap-4">
                                            <div>
                                                <label className="flex items-center space-x-2 cursor-pointer">
                                                    <input
                                                        type="radio"
                                                        name={conflict.id}
                                                        value="source"
                                                        checked={resolutions[conflict.id] === 'source'}
                                                        onChange={() => handleResolution(conflict.id, 'source')}
                                                        className="text-green-600 focus:ring-green-500"
                                                    />
                                                    <div>
                                                        <div className="font-medium text-figma-black">Keep from source branch</div>
                                                        <div className="text-sm text-gray-600">"{conflict.source_value}"</div>
                                                    </div>
                                                </label>
                                            </div>

                                            <div>
                                                <label className="flex items-center space-x-2 cursor-pointer">
                                                    <input
                                                        type="radio"
                                                        name={conflict.id}
                                                        value="target"
                                                        checked={resolutions[conflict.id] === 'target'}
                                                        onChange={() => handleResolution(conflict.id, 'target')}
                                                        className="text-green-600 focus:ring-green-500"
                                                    />
                                                    <div>
                                                        <div className="font-medium text-figma-black">Keep current value</div>
                                                        <div className="text-sm text-gray-600">"{conflict.target_value}"</div>
                                                    </div>
                                                </label>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>

                            <div className="flex justify-end gap-2 mt-4">
                                <button
                                    onClick={() => setStep('select')}
                                    className="px-4 py-2 text-gray-600 hover:text-gray-800 font-medium"
                                >
                                    Back
                                </button>
                                <button
                                    onClick={handleMerge}
                                    disabled={!allConflictsResolved || loading}
                                    className="px-4 py-2 bg-green-600 text-figma-white font-medium rounded-lg hover:bg-green-700 disabled:opacity-50"
                                >
                                    {loading ? 'Merging...' : 'Execute Merge'}
                                </button>
                            </div>
                        </div>
                    )}

                    {step === 'no-conflicts' && (
                        <div>
                            <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg">
                                <div className="flex items-center">
                                    <div className="text-green-600 mr-2 text-xl">✓</div>
                                    <div>
                                        <h3 className="font-medium text-green-800">No conflicts found!</h3>
                                        <p className="text-sm text-green-700 mt-1">
                                            The selected branch can be merged cleanly into <strong>{currSheet?.curr_branch.name}</strong> without any conflicts.
                                        </p>
                                    </div>
                                </div>
                            </div>

                            <div className="flex justify-end gap-2 mt-4">
                                <button
                                    onClick={() => setStep('select')}
                                    className="px-4 py-2 text-gray-600 hover:text-gray-800 font-medium"
                                >
                                    Back
                                </button>
                                <button
                                    onClick={handleMerge}
                                    disabled={loading}
                                    className="px-4 py-2 bg-green-600 text-figma-white font-medium rounded-lg hover:bg-green-700 disabled:opacity-50"
                                >
                                    {loading ? 'Merging...' : 'Proceed with Merge'}
                                </button>
                            </div>
                        </div>
                    )}

                    {step === 'merging' && (
                        <div className="text-center">
                            <div className="mb-4">
                                <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-green-600"></div>
                            </div>
                            <p className="text-figma-black font-medium">Merging branches...</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
