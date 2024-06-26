// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function AddStoryModule($module: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3262065790, $module) as any;
    return $resultPromise;
}

export function CreateTag(tagName: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2945885466, tagName) as any;
    return $resultPromise;
}

export function EditStoryModuleConfig($module: string, config: string, value: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3576513881, $module, config, value) as any;
    return $resultPromise;
}

export function GetOpenStory(): Promise<$models.Story | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2563610799) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function GetStory(projectId: string, refresh: boolean): Promise<$models.Story | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(180270295, projectId, refresh) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function GetStoryImages(): Promise<$models.ImageFile[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(619858253) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType3($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function GetStoryModules(unusedModulesOnly: boolean): Promise<string[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3606125984, unusedModulesOnly) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType4($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function NewStory(projectDirectory: string): Promise<$models.Story | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2808099189, projectDirectory) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function SaveStory(): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3515724504) as any;
    return $resultPromise;
}

export function SetStoryDescription(description: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3596101199, description) as any;
    return $resultPromise;
}

export function SetStoryTitle(name: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2387590379, name) as any;
    return $resultPromise;
}

// Private type creation functions
const $$createType0 = $models.Story.createFrom;
const $$createType1 = $Create.Nullable($$createType0);
const $$createType2 = $models.ImageFile.createFrom;
const $$createType3 = $Create.Array($$createType2);
const $$createType4 = $Create.Array($Create.Any);
