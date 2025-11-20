package combine;

import consts.ConstFont;
import consts.ConstNpc;
import item.Item;
import player.Player;
import player.Service.InventoryService;
import services.Service;
import utils.Util;

public class GiamDinhSach {

    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() != 2) {
            Service.gI().sendDialogMessage(player, "Cần Sách Tuyệt Kỹ và bùa giám định.");
            return;
        }

        Item sachTuyetKy = null;
        Item buaGiamDinh = null;

        for (Item item : player.combineNew.itemsCombine) {
            if (item.isSachTuyetKy() || item.isSachTuyetKy2()) {
                sachTuyetKy = item;
            } else if (item.template.id == 1284) {
                buaGiamDinh = item;
            }
        }

        if (sachTuyetKy == null || buaGiamDinh == null) {
            Service.gI().sendDialogMessage(player, "Cần Sách Tuyệt Kỹ và bùa giám định.");
            return;
        }

        StringBuilder text = new StringBuilder();
        text.append(ConstFont.BOLD_GREEN)
            .append("Giám định ")
            .append(sachTuyetKy.template.name)
            .append(" ?\n")
            .append(ConstFont.BOLD_BLUE)
            .append("Bùa giám định ")
            .append(buaGiamDinh.quantity)
            .append("/1");

        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, text.toString(), "Giám định", "Từ chối");
    }

    public static void giamDinhSach(Player player) {
        if (player.combineNew.itemsCombine.size() != 2) {
            return;
        }

        Item sachTuyetKy = null;
        Item buaGiamDinh = null;

        for (Item item : player.combineNew.itemsCombine) {
            if (item.isSachTuyetKy() || item.isSachTuyetKy2()) {
                sachTuyetKy = item;
            } else if (item.template.id == 1284) {
                buaGiamDinh = item;
            }
        }

        if (sachTuyetKy == null || buaGiamDinh == null || !sachTuyetKy.isHaveOption(217)) {
            Service.gI().sendServerMessage(player, "Còn cái nịt mà giám");
            return;
        }

        int[] options = {77, 103, 50, 108, 94, 14, 80, 81, 175, 5, 214, 216};

        for (int i = 0; i < sachTuyetKy.itemOptions.size(); i++) {
            Item.ItemOption io = sachTuyetKy.itemOptions.get(i);
            if (io.optionTemplate.id == 217) {
                int randomOption = options[Util.nextInt(options.length)];
                int randomValue = Util.nextInt(1, 10 / Util.nextInt(1, 3));
                sachTuyetKy.itemOptions.set(i, new Item.ItemOption(randomOption, randomValue));
            }
        }

        CombineService.gI().sendEffectSuccessCombine(player);
        InventoryService.gI().subQuantityItemsBag(player, buaGiamDinh, 1);
        InventoryService.gI().sendItemBags(player);
        CombineService.gI().reOpenItemCombine(player);
    }
}