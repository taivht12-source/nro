package combine;

/*
 * @Author Coder: Nguyễn Tấn Tài
 * @Description: Ngọc Rồng Kiwi - Máy Chủ Chuẩn Teamobi 2025
 * @Group Zalo: https://zalo.me/g/toiyeuvietnam2025
 */
import consts.ConstNpc;
import item.Item;
import static combine.CombineService.MAX_STAR_ITEM;
import player.Player;
import server.ServerNotify;
import services.ChatGlobalService;
import player.Service.InventoryService;
import services.Service;
import services.TaskService;
import utils.Util;

public class PhaLeHoaTrangBi {

    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() != 1) {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, "Hãy chọn 1 vật phẩm để pha lê hóa", "Đóng");
            return;
        }

        Item item = player.combineNew.itemsCombine.get(0);
        if (!CombineSystem.isTrangBiPhaLeHoa(item)) {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, "Vật phẩm này không thể đục lỗ", "Đóng");
            return;
        }

        int star = 0;
        int epStar = -1; // Giá trị mặc định là chưa ép
        for (Item.ItemOption io : item.itemOptions) {
            if (io.optionTemplate.id == 107) {
                star = io.param; // Lấy số sao hiện tại
            }
            if (io.optionTemplate.id == 102) {
                epStar = io.param; // Lấy số sao đã ép
            }
        }

        // Kiểm tra số sao hiện tại có vượt quá giới hạn tối đa
        if (star >= CombineService.MAX_STAR_ITEM) {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, "Vật phẩm đã đạt tối đa sao pha lê", "Đóng");
            return;
        }

        // Nếu sao < 7, không cần kiểm tra ép
        if (star < 7) {
            processPhaLeHoa(player, item, star);
            return;
        }

        // Nếu sao >= 7, kiểm tra điều kiện ép
        if (epStar == -1) {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, "Vật phẩm chưa được ép, không thể pha lê hóa", "Đóng");
            return;
        }

        if (epStar != star) {
            String message = (star == 8)
                    ? "Cần cường hóa và ép lỗ thứ 8 để tiếp tục pha lê hóa"
                    : "Chưa ép hết lỗ sao, không thể pha lê hóa";
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, message, "Đóng");
            return;
        }

        // Xử lý pha lê hóa khi tất cả điều kiện đã đạt
        processPhaLeHoa(player, item, star);
    }

    private static void processPhaLeHoa(Player player, Item item, int star) {
        player.combineNew.goldCombine = CombineSystem.getGoldPhaLeHoa(star);
        player.combineNew.gemCombine = CombineSystem.getGemPhaLeHoa(star);
        player.combineNew.ratioCombine = CombineSystem.getRatioPhaLeHoa(star);

        String npcSay = item.template.name + "\n|2|";
        for (Item.ItemOption io : item.itemOptions) {
            if (io.optionTemplate.id != 102) {
                npcSay += io.getOptionString() + "\n";
            }
        }
        npcSay += "|7|Tỉ lệ thành công: " + player.combineNew.ratioCombine + "%" + "\n";

        if (player.combineNew.goldCombine <= player.inventory.gold) {
            npcSay += "|1|Cần " + Util.numberToMoney(player.combineNew.goldCombine) + " vàng";
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay,
                    "Nâng cấp\ncần " + player.combineNew.gemCombine + " ngọc", "Nâng cấp 10 lần", "Nâng cấp 100 lần");
        } else {
            npcSay += "Còn thiếu " + Util.numberToMoney(player.combineNew.goldCombine - player.inventory.gold) + " vàng";
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, npcSay, "Đóng");
        }
    }

    public static void phaLeHoa(Player player, int... numm) {
        if (player.idMark != null && !Util.canDoWithTime(player.idMark.getLastTimeCombine(), 500)) {
            return;
        }
        player.idMark.setLastTimeCombine(System.currentTimeMillis());
        int n = 1;
        if (numm.length > 0) {
            n = numm[0];
        }
        if (!player.combineNew.itemsCombine.isEmpty()) {
            int gold = player.combineNew.goldCombine;
            int gem = player.combineNew.gemCombine;
            if (player.inventory.gold < gold) {
                Service.gI().sendThongBao(player, "Không đủ vàng để thực hiện");
                return;
            } else if (player.inventory.gem < gem) {
                Service.gI().sendThongBao(player, "Không đủ ngọc để thực hiện");
                return;
            }
            int num = 0;
            int star = 0;
            boolean success = false;
            int fail = 0;
            Item item = null;
            Item.ItemOption optionStar = null;

            for (int i = 0; i < n; i++) {
                num = i;
                gold = player.combineNew.goldCombine;
                gem = player.combineNew.gemCombine;
                if (player.inventory.gem < gem || player.inventory.gold < gold) {
                    break;
                }

                item = player.combineNew.itemsCombine.get(0);
                if (CombineSystem.isTrangBiPhaLeHoa(item)) {
                    star = 0;
                    optionStar = null;
                    for (Item.ItemOption io : item.itemOptions) {
                        if (io.optionTemplate.id == 107) {
                            star = io.param;
                            optionStar = io;
                            break;
                        }
                    }
                    if (star < MAX_STAR_ITEM) {
                        player.combineNew.goldCombine = CombineSystem.getGoldPhaLeHoa(star);
                        player.combineNew.gemCombine = CombineSystem.getGemPhaLeHoa(star);
                        player.combineNew.ratioCombine = CombineSystem.getRatioPhaLeHoa(star);
                        player.inventory.gold -= gold;
                        player.inventory.gem -= gem;

                        int ratio = 1;
                        boolean succ = true;
                        if (optionStar != null) {
                            switch (optionStar.param) {

                                case 4 ->
                                    ratio *= 120/100;
                                case 5 ->
                                    ratio *= 120/100;
                                case 6 ->
                                    ratio *= 120/100;
                                case 7 ->
                                    ratio *= 120/100;
                                case 8 ->
                                    ratio *= 120/100;
                                case 9 ->
                                    ratio *= 120/100;

                            }

                        }
                        if (Util.isTrue(player.combineNew.ratioCombine, 100 * ratio) && succ) {
                            success = true;
                            break;
                        } else {
                            fail++;
                        }
                    }
                } else {
                    break;
                }
            }

            if (success) {
                star++;
                if (item != null) {
                    if (optionStar == null) {
                        item.itemOptions.add(new Item.ItemOption(107, star));
                    } else {
                        optionStar.param = star;
                    }
                    if (optionStar != null && optionStar.param >= 7) {

                        ChatGlobalService.gI().ThongBaoDapDo(player,"Chúc mừng " + player.name + " vừa pha lê hóa " + "thành công " + item.template.name + " lên " + star + " sao pha lê");
                    }
                }
                if (n > 1 && num > 1) {
                    Service.gI().sendThongBao(player, "Pha lê hóa trang bị lên " + star + " sao thành công, sau " + num + " lần nâng cấp!");
                }
                CombineService.gI().sendEffectSuccessCombine(player);
            } else {
                CombineService.gI().sendEffectFailCombine(player);
            }
            InventoryService.gI().sendItemBags(player);
            Service.gI().sendMoney(player);
            CombineService.gI().reOpenItemCombine(player);
        }

    }
}
