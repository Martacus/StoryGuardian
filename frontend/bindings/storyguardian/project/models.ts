// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Create as $Create} from "@wailsio/runtime";

export class ApplicationConfig {
    "projects": { [_: string]: ProjectDetails };

    /** Creates a new ApplicationConfig instance. */
    constructor($$source: Partial<ApplicationConfig> = {}) {
        if (!("projects" in $$source)) {
            this["projects"] = {};
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new ApplicationConfig instance from a string or object.
     */
    static createFrom($$source: any = {}): ApplicationConfig {
        const $$createField0_0 = $$createType1;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("projects" in $$parsedSource) {
            $$parsedSource["projects"] = $$createField0_0($$parsedSource["projects"]);
        }
        return new ApplicationConfig($$parsedSource as Partial<ApplicationConfig>);
    }
}

export class Entity {
    "Id": string;
    "Name": string;
    "ProjectId": string;
    "Description": string;
    "CreatedAt": string;
    "Tags": string[];

    /** Creates a new Entity instance. */
    constructor($$source: Partial<Entity> = {}) {
        if (!("Id" in $$source)) {
            this["Id"] = "";
        }
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("ProjectId" in $$source)) {
            this["ProjectId"] = "";
        }
        if (!("Description" in $$source)) {
            this["Description"] = "";
        }
        if (!("CreatedAt" in $$source)) {
            this["CreatedAt"] = "";
        }
        if (!("Tags" in $$source)) {
            this["Tags"] = [];
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Entity instance from a string or object.
     */
    static createFrom($$source: any = {}): Entity {
        const $$createField5_0 = $$createType2;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("Tags" in $$parsedSource) {
            $$parsedSource["Tags"] = $$createField5_0($$parsedSource["Tags"]);
        }
        return new Entity($$parsedSource as Partial<Entity>);
    }
}

export class Project {
    "id": string;
    "name": string;
    "location": string;
    "description": string;
    "entities": Entity[];

    /** Creates a new Project instance. */
    constructor($$source: Partial<Project> = {}) {
        if (!("id" in $$source)) {
            this["id"] = "";
        }
        if (!("name" in $$source)) {
            this["name"] = "";
        }
        if (!("location" in $$source)) {
            this["location"] = "";
        }
        if (!("description" in $$source)) {
            this["description"] = "";
        }
        if (!("entities" in $$source)) {
            this["entities"] = [];
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Project instance from a string or object.
     */
    static createFrom($$source: any = {}): Project {
        const $$createField4_0 = $$createType4;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("entities" in $$parsedSource) {
            $$parsedSource["entities"] = $$createField4_0($$parsedSource["entities"]);
        }
        return new Project($$parsedSource as Partial<Project>);
    }
}

export class ProjectDetails {
    "id": string;
    "name": string;
    "location": string;

    /** Creates a new ProjectDetails instance. */
    constructor($$source: Partial<ProjectDetails> = {}) {
        if (!("id" in $$source)) {
            this["id"] = "";
        }
        if (!("name" in $$source)) {
            this["name"] = "";
        }
        if (!("location" in $$source)) {
            this["location"] = "";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new ProjectDetails instance from a string or object.
     */
    static createFrom($$source: any = {}): ProjectDetails {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new ProjectDetails($$parsedSource as Partial<ProjectDetails>);
    }
}

// Private type creation functions
const $$createType0 = ProjectDetails.createFrom;
const $$createType1 = $Create.Map($Create.Any, $$createType0);
const $$createType2 = $Create.Array($Create.Any);
const $$createType3 = Entity.createFrom;
const $$createType4 = $Create.Array($$createType3);
