package presenters

import (
	"goilerplate/domain/models"
	up "goilerplate/usecase/presenters"
)

type groupPresenter struct{}

func NewGroupPresenter() up.IGroupPresenter {
	return &groupPresenter{}
}

func (p *groupPresenter) PresentAll(g *[]models.Group) []up.GroupPresent {
	groups := make([]up.GroupPresent, len(*g))
	for i, gr := range *g {
		groups[i] = up.GroupPresent{ID: gr.ID, Name: gr.Name, Description: gr.Description}
	}

	return groups
}

func (p *groupPresenter) PresentAlongSelectedRoles(g *models.Group, r *[]up.RolePresent) up.GroupAlongSelectedRolePresent {
	group := up.GroupAlongSelectedRolePresent{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Roles:       make([]up.SelectedRolePresent, len(*r)),
	}

	for i, rl := range *r {
		role := up.SelectedRolePresent{
			ID:     rl.ID,
			FaName: rl.FaName,
			EnName: rl.EnName,
		}
		for _, rle := range g.Roles {
			if rle.ID == rl.ID {
				role.Selected = true
			}
		}
		group.Roles[i] = role
	}

	return group
}

func (p *groupPresenter) PresentAlongRolesAndUsers(g *models.Group) up.GroupAlongRolesAndUsersPresent {
	group := up.GroupAlongRolesAndUsersPresent{Name: g.Name, Description: g.Description}

	group.Users = make([]up.UserPresent, len(g.Users))
	for i, u := range g.Users {
		group.Users[i] = up.UserPresent{
			FullName:     u.FirstName + " " + u.LastName,
			MobileNumber: u.MobileNumber,
		}
	}

	group.Roles = make([]up.RolePresent, len(g.Roles))
	for i, r := range g.Roles {
		group.Roles[i] = up.RolePresent{
			FaName: r.FaName,
			EnName: r.EnName,
		}
	}

	return group
}

func (p *groupPresenter) PresentSave(g *models.Group) up.GroupPresent {
	return PresentGroup(g)
}

func (p *groupPresenter) PresentEdit(g *models.Group) up.GroupPresent {
	return PresentGroup(g)
}

func PresentGroup(g *models.Group) up.GroupPresent {
	return up.GroupPresent{ID: g.ID, Name: g.Name, Description: g.Description}
}
