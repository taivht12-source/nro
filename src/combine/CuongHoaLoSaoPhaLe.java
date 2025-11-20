package combine;

/*
 * @Author Coder: Nguyễn Tấn Tài
 * @Description: Ngọc Rồng Kiwi - Máy Chủ Chuẩn Teamobi 2025
 * @Group Zalo: https://zalo.me/g/toiyeuvietnam2025
 */
import consts.ConstNpc;
import item.Item;
import item.Item.ItemOption;
import player.Player;
import player.Service.InventoryService;
import services.Service;
import utils.Util;

public class CuongHoaLoSaoPhaLe {

    private static final int COST = 500000000;

    public static void showInfoCombine(Player player) {
        if (InventoryService.gI().getCountEmptyBag(player) > 0) {
            if (player.combineNew.itemsCombine.size() == 3) {
                Item item = null, Hematite = null, DuiDuc = null;

                for (Item i : player.combineNew.itemsCombine) {
                    if (CombineSystem.isTrangBiPhaLeHoa(i)) {
                        item = i;
                    } else if (i.template.id == 1423) { // Hematite
                        Hematite = i;
                    } else if (i.template.id == 1438) { // Dùi Đục
                        DuiDuc = i;
                    }
                }

                if (item != null && Hematite != null && DuiDuc != null
                        && Hematite.quantity >= 1 && DuiDuc.quantity >= 1) {

                    String npcSay = item.template.name + "\n|2|";
                    for (ItemOption io : Hematite.itemOptions) {
                        npcSay += io.getOptionString() + "\n";
                    }
                    npcSay += "Cường hóa\n" + " Ô sao pha lê thứ 8\n" + item.template.name
                            + "\nTỉ lệ thành công: 50%\n"
                            + "|7| Cần 1 " + Hematite.template.name + "\n|7| Cần 1 " + DuiDuc.template.name + "\nCần "
                            + Util.numberToMoney(COST) + " vàng";

                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay, "Cường Hóa", "Từ chối");
                } else {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, "Bạn chưa bỏ đủ vật phẩm !!!", "Đóng");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, "Cần bỏ đủ vật phẩm yêu cầu", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, "Hành trang cần ít nhất 1 chỗ trống", "Đóng");
        }
    }

    public static void cuongHoaLoSaoPhaLe(Player player) {
        if (InventoryService.gI().getCountEmptyBag(player) > 0) {
            if (player.inventory.gold < COST) {
                Service.gI().sendThongBao(player, "Con cần thêm vàng để cường hóa...");
                return;
            }

            if (!player.combineNew.itemsCombine.isEmpty()) {
                Item item = null;
                Item Hematite = null;
                Item DuiDuc = null;

                for (Item i : player.combineNew.itemsCombine) {
                    if (CombineSystem.isTrangBiPhaLeHoa(i)) {
                        item = i;
                    } else if (i.template.id == 1423) { // ID của Hematite
                        Hematite = i;
                    } else if (i.template.id == 1438) { // ID của DuiDuc
                        DuiDuc = i;
                    }
                }

                if (item != null && Hematite != null && DuiDuc != null
                        && Hematite.quantity >= 1 && DuiDuc.quantity >= 1) {

                    int star = 0;
                    ItemOption optionStar = null;

                    for (ItemOption io : item.itemOptions) {
                        if (io.optionTemplate.id == 107) {
                            star = io.param;
                        }
                        if (io.optionTemplate.id == 228) {
                            optionStar = io;
                        }
                    }

                    if (star == 8 && optionStar == null) {
                        item.itemOptions.add(new ItemOption(218, 0));
                        item.itemOptions.add(new ItemOption(228, 8));

                        player.inventory.gold -= COST;
                        if (Util.isTrue(50, 100)) {
                            CombineService.gI().sendEffectSuccessCombine(player);
                        } else {
                            CombineService.gI().sendEffectFailCombine(player);
                        }
                    } else if (star == 9 && optionStar != null && optionStar.param == 8) {
                        player.inventory.gold -= COST;

                        if (Util.isTrue(50, 100)) {
                            optionStar.param += 1;  // Thêm 1 vào param của optionStar
                            CombineService.gI().sendEffectSuccessCombine(player);

                            CombineService.gI().sendEffectSuccessCombine(player);

                            Service.gI().sendThongBao(player, "Trang bị của bạn đã cường hóa thành công lên sao thứ 9!");
                        } else {
                            CombineService.gI().sendEffectFailCombine(player);
                        }
                    } else if (optionStar != null && optionStar.param >= 9) {
                        Service.gI().sendThongBao(player, "Trang bị của bạn đã đạt tối đa, không thể cường hóa thêm.");
                        return;
                    } else {
                        Service.gI().sendThongBao(player, "Cường hóa không hợp lệ, vui lòng nâng cấp trang bị lên 9 sao.");
                        return;
                    }

                    InventoryService.gI().subQuantityItemsBag(player, Hematite, 1);
                    InventoryService.gI().subQuantityItemsBag(player, DuiDuc, 1);
                    Service.gI().sendMoney(player);
                    InventoryService.gI().sendItemBags(player);
                    CombineService.gI().reOpenItemCombine(player);
                    CombineService.gI().sendEffectCombineDB(player, item.template.iconID);

                } else {
                    Service.gI().sendThongBao(player, "Vật phẩm không hợp lệ hoặc không đủ số lượng.");
                }
            }
        }
    }
}
