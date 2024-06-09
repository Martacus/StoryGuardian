// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Create as $Create} from "@wailsio/runtime";

export class Project {
    "name": string;
    "location": string;

    /** Creates a new Project instance. */
    constructor($$source: Partial<Project> = {}) {
        if (!("name" in $$source)) {
            this["name"] = "";
        }
        if (!("location" in $$source)) {
            this["location"] = "";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Project instance from a string or object.
     */
    static createFrom($$source: any = {}): Project {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new Project($$parsedSource as Partial<Project>);
    }
}

export class ProjectConfig {
    "projects": { [_: string]: Project };

    /** Creates a new ProjectConfig instance. */
    constructor($$source: Partial<ProjectConfig> = {}) {
        if (!("projects" in $$source)) {
            this["projects"] = {};
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new ProjectConfig instance from a string or object.
     */
    static createFrom($$source: any = {}): ProjectConfig {
        const $$createField0_0 = $$createType1;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("projects" in $$parsedSource) {
            $$parsedSource["projects"] = $$createField0_0($$parsedSource["projects"]);
        }
        return new ProjectConfig($$parsedSource as Partial<ProjectConfig>);
    }
}

// Private type creation functions
const $$createType0 = Project.createFrom;
const $$createType1 = $Create.Map($Create.Any, $$createType0);
