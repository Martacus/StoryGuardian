using System;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.UI;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using StoryGuardian.Data;

[assembly: HostingStartup(typeof(StoryGuardian.Areas.Identity.IdentityHostingStartup))]
namespace StoryGuardian.Areas.Identity
{
    public class IdentityHostingStartup : IHostingStartup
    {
        public void Configure(IWebHostBuilder builder)
        {
            builder.ConfigureServices((context, services) => {
                services.AddDbContext<StoryGuardianDbContext>(options =>
                    options.UseSqlServer(
                        context.Configuration.GetConnectionString("StoryGuardianDbContextConnection")));

                services.AddDefaultIdentity<IdentityUser>(options => options.SignIn.RequireConfirmedAccount = false)
                    .AddRoles<IdentityRole>().AddEntityFrameworkStores<StoryGuardianDbContext>();
            });
        }
    }
}