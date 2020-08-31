using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Threading.Tasks;

namespace StoryGuardian.Models.Story
{
    public class EntityModel
    {
        public Guid Id { get; set; }

        [Required, StringLength(50)]
        public string Name { get; set; }

        public string Description { get; set; }

        public string ImageLocation { get; set; }

        public StoryModel Story {get; set;}
    }
}
