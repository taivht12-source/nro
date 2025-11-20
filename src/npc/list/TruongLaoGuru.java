package npc.list;
import clan.Clan;
import consts.ConstNpc;
import consts.ConstPlayer;
import item.Item;
import java.util.ArrayList;
import npc.Npc;
import player.Player;
import services.ItemService;
import map.Service.NpcService;
import services.RewardService;
import services.Service;
import shop.ShopService;
import services.TaskService;
import map.Service.ChangeMapService;
import services.func.Input;
import player.Service.PlayerService;
import skill.Skill;
import utils.Logger;
import utils.SkillUtil;
import utils.TimeUtil;
import utils.Util;

public class TruongLaoGuru extends Npc {

    public TruongLaoGuru(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            if (!TaskService.gI().checkDoneTaskTalkNpc(player, this)) {
                if (player.gender != ConstPlayer.NAMEC) {
                    NpcService.gI().createTutorial(player, tempId, avartar, "Con hãy về hành tinh của mình mà thể hiện");
                    return;
                }
                ArrayList<String> menu = new ArrayList<>();
                if (!player.canReward) {
                    menu.add("Nhiệm vụ");
                    menu.add("Học\nKỹ năng");
                    Clan clan = player.clan;
                    if (clan != null) {
                        menu.add("Về khu\nvực bang");
                        if (clan.isLeader(player)) {
                            menu.add("Giải tán\nBang hội");
                        }
                    }
                } else {
                    menu.add("Giao\nLân con");
                }
                String[] menus = menu.toArray(String[]::new);
                createOtherMenu(player, ConstNpc.BASE_MENU,
                        "Chào con, ta rất vui khi gặp được con\nCon muốn làm gì nào ?", menus);
            }
        }
    }

    
    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            if (player.canReward) {
                RewardService.gI().rewardLancon(player);
                return;
            }

            if (player.idMark.isBaseMenu()) {
                handleBaseMenu(player, select);
            } else if (player.idMark.getIndexMenu() == 3) {
                handleClanDissolution(player, select);
            } else if (player.idMark.getIndexMenu() == 12) {
                handleSkillLearning(player, select);
            }
        }
    }

    private void handleBaseMenu(Player player, int select) {
        switch (select) {
            case 0 -> NpcService.gI().createTutorial(player, tempId, avartar, player.playerTask.taskMain.subTasks.get(player.playerTask.taskMain.index).name);
            case 1 -> handleSkillLearningMenu(player);
            case 2 -> handleClanMenu(player);
            case 3 -> handleClanLeaderMenu(player);
        }
    }

    private void handleSkillLearningMenu(Player player) {
        if (player.LearnSkill.Time != -1) {
            var ngoc = 5;
            var time = player.LearnSkill.Time - System.currentTimeMillis();
            if (time / 600_000 >= 2) {
                ngoc += time / 600_000;
            }
            String[] subName = ItemService.gI().getTemplate(player.LearnSkill.ItemTemplateSkillId).name.split("");
            byte level = Byte.parseByte(subName[subName.length - 1]);
            createOtherMenu(player, 12,
                    "Con đang học kỹ năng\n" + SkillUtil.findSkillTemplate(SkillUtil.getTempSkillSkillByItemID(player.LearnSkill.ItemTemplateSkillId)).name
                            + " cấp " + level + "\nThời gian còn lại " + TimeUtil.getTime(time),
                    "Học Cấp tốc " + ngoc + " ngọc", "Huỷ", "Bỏ qua");
        } else {
            ShopService.gI().opendShop(player, "QUY_LAO", false);
        }
    }

    private void handleClanMenu(Player player) {
        if (player.clan != null) {
            ChangeMapService.gI().changeMapNonSpaceship(player, 153, Util.nextInt(100, 200), 432);
        }
    }

    private void handleClanLeaderMenu(Player player) {
        if (player.clan != null && player.clan.isLeader(player)) {
            createOtherMenu(player, 3, "Con có chắc muốn giải tán bang hội không?", "Đồng ý", "Từ chối");
        }
    }

    private void handleClanDissolution(Player player, int select) {
        if (player.clan != null && player.clan.isLeader(player) && select == 0) {
            Input.gI().createFormGiaiTanBangHoi(player);
        }
    }

    private void handleSkillLearning(Player player, int select) {
        switch (select) {
            case 0 -> {
                var time = player.LearnSkill.Time - System.currentTimeMillis();
                var ngoc = 5;
                if (time / 600_000 >= 2) {
                    ngoc += time / 600_000;
                }
                if (player.inventory.gem < ngoc) {
                    Service.gI().sendThongBao(player, "Bạn không có đủ ngọc");
                    return;
                }
                player.inventory.subGem(ngoc);
                player.LearnSkill.Time = -1;
                learnSkill(player);
            }
            case 1 -> createOtherMenu(player, 13, "Con có muốn huỷ học kỹ năng này và nhận lại 50% số tiềm năng không ?", "Ok", "Đóng");
        }
    }

    private void learnSkill(Player player) {
        try {
            String[] subName = ItemService.gI().getTemplate(player.LearnSkill.ItemTemplateSkillId).name.split("");
            byte level = Byte.parseByte(subName[subName.length - 1]);
            Skill curSkill = SkillUtil.getSkillByItemID(player, player.LearnSkill.ItemTemplateSkillId);

            if (curSkill.point == 0) {
                player.BoughtSkill.add((int) player.LearnSkill.ItemTemplateSkillId);
                curSkill = SkillUtil.createSkill(SkillUtil.getTempSkillSkillByItemID(player.LearnSkill.ItemTemplateSkillId), level);
                SkillUtil.setSkill(player, curSkill);
                var msg = Service.gI().messageSubCommand((byte) 62);
                msg.writer().writeShort(curSkill.skillId);
                player.sendMessage(msg);
                msg.cleanup();
                PlayerService.gI().sendInfoHpMpMoney(player);
            } else {
                curSkill = SkillUtil.createSkill(SkillUtil.getTempSkillSkillByItemID(player.LearnSkill.ItemTemplateSkillId), level);
                player.BoughtSkill.add((int) player.LearnSkill.ItemTemplateSkillId);
                SkillUtil.setSkill(player, curSkill);
                var msg = Service.gI().messageSubCommand((byte) 62);
                msg.writer().writeShort(curSkill.skillId);
                player.sendMessage(msg);
                msg.cleanup();
                PlayerService.gI().sendInfoHpMpMoney(player);
            }
        } catch (Exception e) {
            Logger.log(e.toString());
        }
    }
}