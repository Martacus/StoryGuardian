// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function CreateTag(tagName: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(798792175, tagName) as any;
    return $resultPromise;
}

export function GetStory(projectId: string): Promise<$models.Story | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(941168768, projectId) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function GetStoryImages(): Promise<$models.ImageFile[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4053663178) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType3($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function NewStory(projectDirectory: string): Promise<$models.Story | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1310826610, projectDirectory) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function SaveStory(): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2504230633) as any;
    return $resultPromise;
}

export function SetStoryDescription(description: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1384865094, description) as any;
    return $resultPromise;
}

export function SetStoryTitle(name: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1996276650, name) as any;
    return $resultPromise;
}

// Private type creation functions
const $$createType0 = $models.Story.createFrom;
const $$createType1 = $Create.Nullable($$createType0);
const $$createType2 = $models.ImageFile.createFrom;
const $$createType3 = $Create.Array($$createType2);
