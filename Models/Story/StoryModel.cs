using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Identity;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using System.Linq;
using System.Threading.Tasks;

namespace StoryGuardian.Models.Story
{
    public class StoryModel
    {
        public Guid Id { get; set; }

        [Required, StringLength(50)]
        public string Name { get; set; }

        public string Description { get; set; }

        public string ImageLocation { get; set; }

        public IdentityUser User { get; set; }

        //[NotMapped]
        //[Display(Name = "Story Image")]
        //public IFormFile FormFile { get; set; }
    }
}
