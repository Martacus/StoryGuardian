// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Create as $Create} from "@wailsio/runtime";

export class ApplicationConfig {
    "projects": { [_: string]: ProjectDetails };
    "openProject": ProjectDetails;

    /** Creates a new ApplicationConfig instance. */
    constructor($$source: Partial<ApplicationConfig> = {}) {
        if (!("projects" in $$source)) {
            this["projects"] = {};
        }
        if (!("openProject" in $$source)) {
            this["openProject"] = (new ProjectDetails());
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new ApplicationConfig instance from a string or object.
     */
    static createFrom($$source: any = {}): ApplicationConfig {
        const $$createField0_0 = $$createType1;
        const $$createField1_0 = $$createType0;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("projects" in $$parsedSource) {
            $$parsedSource["projects"] = $$createField0_0($$parsedSource["projects"]);
        }
        if ("openProject" in $$parsedSource) {
            $$parsedSource["openProject"] = $$createField1_0($$parsedSource["openProject"]);
        }
        return new ApplicationConfig($$parsedSource as Partial<ApplicationConfig>);
    }
}

export class Entity {
    "id": string;
    "name": string;
    "storyId": string;
    "description": string;
    "createdAt": string;
    "tags": string[];
    "relations": string[];

    /** Creates a new Entity instance. */
    constructor($$source: Partial<Entity> = {}) {
        if (!("id" in $$source)) {
            this["id"] = "";
        }
        if (!("name" in $$source)) {
            this["name"] = "";
        }
        if (!("storyId" in $$source)) {
            this["storyId"] = "";
        }
        if (!("description" in $$source)) {
            this["description"] = "";
        }
        if (!("createdAt" in $$source)) {
            this["createdAt"] = "";
        }
        if (!("tags" in $$source)) {
            this["tags"] = [];
        }
        if (!("relations" in $$source)) {
            this["relations"] = [];
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Entity instance from a string or object.
     */
    static createFrom($$source: any = {}): Entity {
        const $$createField5_0 = $$createType2;
        const $$createField6_0 = $$createType2;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("tags" in $$parsedSource) {
            $$parsedSource["tags"] = $$createField5_0($$parsedSource["tags"]);
        }
        if ("relations" in $$parsedSource) {
            $$parsedSource["relations"] = $$createField6_0($$parsedSource["relations"]);
        }
        return new Entity($$parsedSource as Partial<Entity>);
    }
}

export class ImageFile {
    "location": string;
    "name": string;

    /** Creates a new ImageFile instance. */
    constructor($$source: Partial<ImageFile> = {}) {
        if (!("location" in $$source)) {
            this["location"] = "";
        }
        if (!("name" in $$source)) {
            this["name"] = "";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new ImageFile instance from a string or object.
     */
    static createFrom($$source: any = {}): ImageFile {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new ImageFile($$parsedSource as Partial<ImageFile>);
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

export class Relation {
    "id": string;
    "name": string;
    "entityOne": string;
    "entityTwo": string;
    "description": string;

    /** Creates a new Relation instance. */
    constructor($$source: Partial<Relation> = {}) {
        if (!("id" in $$source)) {
            this["id"] = "";
        }
        if (!("name" in $$source)) {
            this["name"] = "";
        }
        if (!("entityOne" in $$source)) {
            this["entityOne"] = "";
        }
        if (!("entityTwo" in $$source)) {
            this["entityTwo"] = "";
        }
        if (!("description" in $$source)) {
            this["description"] = "";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Relation instance from a string or object.
     */
    static createFrom($$source: any = {}): Relation {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new Relation($$parsedSource as Partial<Relation>);
    }
}

export class RelationInfo {
    "id": string;
    "name": string;
    "entityOne": string;
    "entityTwo": string;
    "description": string;
    "toName": string;

    /** Creates a new RelationInfo instance. */
    constructor($$source: Partial<RelationInfo> = {}) {
        if (!("id" in $$source)) {
            this["id"] = "";
        }
        if (!("name" in $$source)) {
            this["name"] = "";
        }
        if (!("entityOne" in $$source)) {
            this["entityOne"] = "";
        }
        if (!("entityTwo" in $$source)) {
            this["entityTwo"] = "";
        }
        if (!("description" in $$source)) {
            this["description"] = "";
        }
        if (!("toName" in $$source)) {
            this["toName"] = "";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new RelationInfo instance from a string or object.
     */
    static createFrom($$source: any = {}): RelationInfo {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new RelationInfo($$parsedSource as Partial<RelationInfo>);
    }
}

export class Story {
    "id": string;
    "name": string;
    "location": string;
    "description": string;
    "entities": Entity[];

    /** Creates a new Story instance. */
    constructor($$source: Partial<Story> = {}) {
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
     * Creates a new Story instance from a string or object.
     */
    static createFrom($$source: any = {}): Story {
        const $$createField4_0 = $$createType4;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("entities" in $$parsedSource) {
            $$parsedSource["entities"] = $$createField4_0($$parsedSource["entities"]);
        }
        return new Story($$parsedSource as Partial<Story>);
    }
}

// Private type creation functions
const $$createType0 = ProjectDetails.createFrom;
const $$createType1 = $Create.Map($Create.Any, $$createType0);
const $$createType2 = $Create.Array($Create.Any);
const $$createType3 = Entity.createFrom;
const $$createType4 = $Create.Array($$createType3);
