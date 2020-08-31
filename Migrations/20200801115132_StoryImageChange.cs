using System;
using Microsoft.EntityFrameworkCore.Migrations;

namespace StoryGuardian.Migrations
{
    public partial class StoryImageChange : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "Image",
                table: "Stories");

            migrationBuilder.AddColumn<string>(
                name: "ImageLocation",
                table: "Stories",
                nullable: true);
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "ImageLocation",
                table: "Stories");

            migrationBuilder.AddColumn<byte[]>(
                name: "Image",
                table: "Stories",
                type: "varbinary(max)",
                nullable: true);
        }
    }
}
