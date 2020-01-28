package environment

import "git.ringcentral.com/archops/goFsync/core/user"

func AddNewEnv(hostname, envName string, ctx *user.GlobalCTX) error {

	hostID := ctx.Config.Hosts[hostname]

	env, err := ApiGet(hostname, envName, ctx)
	if err != nil {
		return err
	}

	repo := DbGetRepo(hostID, ctx)
	DbInsert(hostID, env.Name, repo, "absent", env.ID, SvnDirInfo{}, ctx)

	return nil
}
