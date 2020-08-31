using Microsoft.AspNetCore.Identity;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace StoryGuardian.Utility
{
    public class RoleHelper
    {
        public RoleHelper()
        {

        }

        public async Task AssignRoles(IdentityUser user, UserManager<IdentityUser> _userManager)
        {
            if (user.Email == "martvdham@gmail.com")
            {
                await _userManager.AddToRoleAsync(user, "Admin");
            }
            else
            {
                await _userManager.AddToRoleAsync(user, "User");
            }
        }

    }
}
