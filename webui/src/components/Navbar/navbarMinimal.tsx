import { ReactComponent as Logo } from "@assets/logo.svg";

import { useLocation, useNavigate } from "react-router-dom";
import { Avatar, Center, createStyles, Navbar, Stack, Tooltip, UnstyledButton } from "@mantine/core";
import { FC, ReactNode } from "react";
import { useRecoilState, useRecoilValue } from "recoil";
import { animated, easings, useSpring } from "react-spring";
import { authState } from "@stores/auth/state";
import { navbarState } from "@src/stores/navbar/state";
import { useAuthActions } from "@src/stores/auth/useAuthActions";

const useStyles = createStyles(theme => ({
  link: {
    width: theme.spacing.lg * 2 + theme.spacing.xs / 2,
    height: theme.spacing.lg * 2 + theme.spacing.xs / 2,
    borderRadius: theme.radius.md,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    color: theme.colors.dark[0],

    "&:hover": {
      backgroundColor: theme.fn.rgba(theme.colors["dg-prim"][1], 0.25),
      "& i": {
        color: theme.colors["dg-prim"][3],
      },
    },

    "& i": {
      color: theme.colors.gray[4],
    },
  },

  active: {
    "&": {
      backgroundColor: theme.fn.rgba(theme.colors["dg-prim"][1], 0.25),
      "& i": {
        color: theme.colors["dg-prim"][3],
      },
    },
  },
}));

interface NavbarLinkProps {
  icon: ReactNode;
  label: string;
  roles?: string[];
  url?: string;

  onClick?(): void;
}

function NavbarLink({ icon, label, onClick, roles, url }: NavbarLinkProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { classes, cx } = useStyles();
  const userInfo = useRecoilValue(authState.userInfo);

  if (url) {
    onClick = () => navigate(`${url}`);
  }

  if (!roles || userInfo?.roles.some(role => roles.includes(role))) {
    return (
      <Tooltip label={label} position="right" withArrow transitionDuration={0}>
        <UnstyledButton
          onClick={onClick}
          className={cx(classes.link, { [classes.active]: url && location.pathname.startsWith(url) })}
        >
          {icon}
        </UnstyledButton>
      </Tooltip>
    );
  }
  return null;
}

const ExtensionToggle: FC<{ canOpen: boolean }> = ({ canOpen }) => {
  const [isOpen, setIsOpen] = useRecoilState(navbarState.open);

  const animStyles = useSpring({
    transform: `rotateY(${isOpen ? 180 : 0}deg)`,
    config: {
      easing: easings.easeInOutQuad,
    },
  });

  if (!canOpen) return null;

  return (
    <NavbarLink
      label={isOpen ? "Collapse Navbar" : "Extend Navbar"}
      icon={<animated.i className={"fas fa-arrow-right"} style={animStyles} />}
      onClick={() => setIsOpen(!isOpen)}
    />
  );
};

// CanOpen indicates if this pages has a NavbarExtension tree
export const NavbarMinimal: FC<{ canOpen: boolean }> = ({ canOpen }) => {
  const navigate = useNavigate();
  const { logoutUser } = useAuthActions();
  const userInfo = useRecoilValue(authState.userInfo);
  const isExtOpen = useRecoilValue(navbarState.open);

  const handleLoginButton = () => {
    navigate("/login", {
      state: {
        from: window.location.pathname,
      },
    });
  };

  return (
    <Navbar width={{ base: 70 }} p="sm" style={isExtOpen ? { borderColor: "transparent" } : undefined}>
      <Center>
        <Logo />
      </Center>
      <Navbar.Section grow mt={30}>
        <Stack align="center" spacing={10}>
          <NavbarLink icon={<i className="fa fa-house-blank" />} onClick={() => navigate("/")} label={"Home"} />
          <NavbarLink
            icon={<i className="fas fa-swords" />}
            label={"Staff"}
            url={"/staff"}
            roles={["developer", "staff"]}
          />
          <NavbarLink
            icon={<i className="fas fa-user-police" />}
            label={"ANG"}
            url={"/police"}
            roles={["developer", "staff", "police"]}
          />
          <NavbarLink
            icon={<i className="fas fa-truck-medical" />}
            label={"Ambulance"}
            url={"/ambu"}
            roles={["developer", "staff", "ambulance"]}
          />
          <NavbarLink icon={<i className="fas fa-code" />} label={"Developer"} url={"/dev"} roles={["developer"]} />
        </Stack>
      </Navbar.Section>
      <Navbar.Section mt={50}>
        <Stack align="center" spacing={10}>
          {userInfo ? (
            <>
              <Avatar src={userInfo.avatarUrl} size={45} />
              <ExtensionToggle canOpen={canOpen} />
              <NavbarLink
                icon={<i className="fas fa-right-from-bracket" />}
                label="Logout"
                onClick={() => logoutUser()}
              />
            </>
          ) : (
            <NavbarLink icon={<i className="fas fa-right-to-bracket" />} label="Login" onClick={handleLoginButton} />
          )}
        </Stack>
      </Navbar.Section>
    </Navbar>
  );
};
