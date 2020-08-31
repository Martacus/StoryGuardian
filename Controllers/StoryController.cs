using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using StoryGuardian.Data;
using StoryGuardian.Models;
using StoryGuardian.Models.Story;

namespace StoryGuardian.Controllers
{
    public class StoryController : Controller
    {
        private readonly StoryGuardianDbContext _context;
        private readonly UserManager<IdentityUser> _userManager;
        private readonly IWebHostEnvironment _hostEnvironment;

        public StoryController(StoryGuardianDbContext context, UserManager<IdentityUser> userManager, IWebHostEnvironment hostEnvironment)
        {
            _context = context;
            _userManager = userManager;
            _hostEnvironment = hostEnvironment;
        }

        public IActionResult Index()
        {
            MainStoryModel model = new MainStoryModel();
            model.Stories = _context.Stories.Where(x => x.User == _userManager.GetUserAsync(User).Result).ToList();
            return View(model);
        }

        public async Task<IActionResult> CreateStory(StoryModel model)
        {

            model.User = await _userManager.GetUserAsync(User);
            if (ModelState.IsValid)
            {
                try
                {
                    //if (model.FormFile.Length > 0)
                    //{
                    //    string wwwRootPath = _hostEnvironment.WebRootPath;
                    //    string fileName = Path.GetFileNameWithoutExtension(model.FormFile.FileName);
                    //    string extension = Path.GetExtension(model.FormFile.FileName);
                    //    fileName = fileName + DateTime.Now.ToString("yyyymmddss") + extension;
                    //    string path = Path.Combine(wwwRootPath + "\\images", fileName);
                    //    string modelPath = Path.Combine("\\images", fileName);

                    //    model.ImageLocation = modelPath;
                    //    using (var fileStream = new FileStream(path, FileMode.Create))
                    //    {
                    //        await model.FormFile.CopyToAsync(fileStream);
                    //    }
                    //}

                    _context.Stories.Add(model);
                    await _context.SaveChangesAsync();
                } catch(Exception e)
                {
                    //todo: exception logging
                }
            }

            return Redirect("/");
        }
    }
}