using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace StoryGuardian.Models.Story
{
    public class EntityCollection
    {
        public string Name { get; set; }
        public string Description { get; set; }

        public StoryModel Story { get; set; }



    }
}
