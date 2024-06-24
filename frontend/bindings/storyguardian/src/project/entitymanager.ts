// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function AddEntityModule(entityId: string, $module: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(32823016, entityId, $module) as any;
    return $resultPromise;
}

export function CreateEntity(entity: $models.Entity): Promise<$models.Entity> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3380210067, entity) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType0($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function CreateRelation(entityId: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3916189964, entityId) as any;
    return $resultPromise;
}

export function GetEntity(entityId: string): Promise<$models.Entity | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2607049459, entityId) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function GetEntityModules(entityId: string, unusedModulesOnly: boolean): Promise<string[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3904918628, entityId, unusedModulesOnly) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType2($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function GetRelation(relationId: string): Promise<$models.Relation> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3307714412, relationId) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType3($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function LoadEntities(): Promise<$models.Entity[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2421700187) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType4($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function LoadRelationInfo(entityId: string, paginationStart: number, amount: number): Promise<$models.RelationInfo[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2199069336, entityId, paginationStart, amount) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType6($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

export function SaveEntity(entityId: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2482894830, entityId) as any;
    return $resultPromise;
}

export function SetEntityDescription(entityId: string, description: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1915442115, entityId, description) as any;
    return $resultPromise;
}

export function SetEntityName(entityId: string, name: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2645553516, entityId, name) as any;
    return $resultPromise;
}

export function SetRelationDescription(relationId: string, description: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1907539802, relationId, description) as any;
    return $resultPromise;
}

export function SetRelationName(relationId: string, name: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2818325651, relationId, name) as any;
    return $resultPromise;
}

// Private type creation functions
const $$createType0 = $models.Entity.createFrom;
const $$createType1 = $Create.Nullable($$createType0);
const $$createType2 = $Create.Array($Create.Any);
const $$createType3 = $models.Relation.createFrom;
const $$createType4 = $Create.Array($$createType0);
const $$createType5 = $models.RelationInfo.createFrom;
const $$createType6 = $Create.Array($$createType5);