import { House, Menu, Moon, Sun, X } from "lucide-react";
import { Link } from "react-router";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "~/components/ui/navigation-menu";
import { useAppContext } from "~/context/useAppContext";
import { Button } from "./ui/button";
import { cn } from "~/lib/utils";
import { useState } from "react";
import SearchCompany from "./SearchCompany";
import GithubSvg from "~/assets/GithubSvg";

export default function Navbar() {
  return (
    <div className="bg-background font-geist sticky top-0 z-50 flex h-(--header-height) w-full">
      <div className="flex w-full items-center justify-between px-2 sm:px-6">
        <div className="hidden h-full items-center justify-between md:flex">
          <NavigationMenu>
            <NavigationMenuList>
              <NavigateLink title="" icon={<House size={20} />} path="/" />
              <NavigateLink title="All Problems" path="/all-problems" />
            </NavigationMenuList>
          </NavigationMenu>
        </div>

        <MobileNav />

        <div className="flex h-full items-center justify-between gap-4">
          <Link
            to="https://github.com/shadil-rayyan/visor-leetcode"
            target="_blank"
            className="hidden sm:flex"
          >
            <Button variant="ghost" size="icon" className="h-8 cursor-pointer">
              <GithubSvg />
            </Button>
          </Link>
          <ToggleTheme />
          <SearchCompany />
        </div>
      </div>
    </div>
  );
}

function MobileNav() {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <>
      <div className="mr-2 flex md:hidden">
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8 p-0 [&>svg]:size-5"
          onClick={() => {
            setMenuOpen((prev) => {
              const next = !prev;
              document.body.style.overflow = next ? "hidden" : "";
              return next;
            });
          }}
        >
          {!menuOpen ? <Menu /> : <X />}
        </Button>
      </div>

      <div
        aria-hidden={!menuOpen}
        className={cn(
          "bg-background/80 fixed inset-x-0 bottom-0 z-40 backdrop-blur-sm",
          "transition-[opacity,transform] duration-200 ease-out",
          menuOpen ? "opacity-100" : "pointer-events-none opacity-0",
        )}
        style={{ top: "var(--header-height)" }}
        onClick={() => {
          setMenuOpen(false);
          document.body.style.overflow = "";
        }}
      >
        <nav className="flex flex-col gap-2 p-4">
          <Link
            to="/"
            className="hover:bg-accent flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium"
            onClick={() => {
              setMenuOpen(false);
              document.body.style.overflow = "";
            }}
          >
            <House size={18} />
            Home
          </Link>
          <Link
            to="/all-problems"
            className="hover:bg-accent flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium"
            onClick={() => {
              setMenuOpen(false);
              document.body.style.overflow = "";
            }}
          >
            All Problems
          </Link>
        </nav>
      </div>
    </>
  );
}

function NavigateLink({
  title,
  path,
  icon,
}: {
  title: string;
  path: string;
  icon?: React.ReactNode;
}) {
  return (
    <NavigationMenuItem>
      <NavigationMenuLink
        asChild
        className={cn(
          navigationMenuTriggerStyle(),
          "h-8 gap-1.5 px-2.5 text-sm",
          !title && "p-2",
        )}
      >
        <Link to={path}>
          {title || icon}
        </Link>
      </NavigationMenuLink>
    </NavigationMenuItem>
  );
}

function ToggleTheme() {
  const { setTheme, theme } = useAppContext();
  return (
    <Button
      variant="ghost"
      size="icon"
      className="h-8 w-8 cursor-pointer p-0"
      onClick={() => setTheme(theme === "light" ? "dark" : "light")}
    >
      {theme === "light" ? <Sun /> : <Moon />}
    </Button>
  );
}
